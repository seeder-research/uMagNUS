package engine64

// Exchange interaction (Heisenberg + Dzyaloshinskii-Moriya)
// See also opencl/exchange.cl and opencl/dmi.cl

import (
	"math"
	"unsafe"

	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
	"github.com/seeder-research/uMagNUS/util"
)

var (
	Aex    = NewScalarParam("Aex", "J/m", "Exchange stiffness", &lex2)
	Dind   = NewScalarParam("Dind", "J/m2", "Interfacial Dzyaloshinskii-Moriya strength", &din2)
	Dbulk  = NewScalarParam("Dbulk", "J/m2", "Bulk Dzyaloshinskii-Moriya strength", &dbulk2)
	lex2   exchParam // inter-cell Aex
	din2   exchParam // inter-cell Dind
	dbulk2 exchParam // inter-cell Dbulk

	B_exch     = NewVectorField("B_exch", "T", "Exchange field", AddExchangeField)
	E_exch     = NewScalarValue("E_exch", "J", "Total exchange energy (including the DMI energy)", GetExchangeEnergy)
	Edens_exch = NewScalarField("Edens_exch", "J/m3", "Total exchange energy density (including the DMI energy density)", AddExchangeEnergyDensity)

	// Average exchange coupling with neighbors. Useful to debug inter-region exchange
	ExchCoupling = NewScalarField("ExchCoupling", "arb.", "Average exchange coupling with neighbors", exchangeDecode)
	DindCoupling = NewScalarField("DindCoupling", "arb.", "Average DMI coupling with neighbors", dindDecode)

	OpenBC = false
)

var AddExchangeEnergyDensity = makeEdensAdder(&B_exch, -0.5) // TODO: normal func

func init() {
	registerEnergy(GetExchangeEnergy, AddExchangeEnergyDensity)
	DeclFunc("ext_ScaleExchange", ScaleInterExchange, "Re-scales exchange coupling between two regions.")
	DeclFunc("ext_InterExchange", InterExchange, "Sets exchange coupling between two regions.")
	DeclFunc("ext_ScaleDind", ScaleInterDind, "Re-scales Dind coupling between two regions.")
	DeclFunc("ext_InterDind", InterDind, "Sets Dind coupling between two regions.")
	DeclVar("OpenBC", &OpenBC, "Use open boundary conditions (default=false)")
	lex2.init(Aex)
	din2.init(Dind)
	dbulk2.init(Dbulk)
}

// Adds the current exchange field to dst
func AddExchangeField(dst *data.Slice) {
	inter := !Dind.isZero()
	bulk := !Dbulk.isZero()
	ms := Msat.MSlice()
	defer ms.Recycle()
	switch {
	case !inter && !bulk:
		opencl.AddExchange(dst, M.Buffer(), lex2.Gpu(), ms, regions.Gpu(), M.Mesh())
	case inter && !bulk:
		Refer("mulkers2017")
		opencl.AddDMI(dst, M.Buffer(), lex2.Gpu(), din2.Gpu(), ms, regions.Gpu(), M.Mesh(), OpenBC) // dmi+exchange
	case bulk && !inter:
		opencl.AddDMIBulk(dst, M.Buffer(), lex2.Gpu(), dbulk2.Gpu(), ms, regions.Gpu(), M.Mesh(), OpenBC) // dmi+exchange
		// TODO: add ScaleInterDbulk and InterDbulk
	case inter && bulk:
		util.Fatal("Cannot have induced and interfacial DMI at the same time")
	}
}

// Set dst to the average exchange coupling per cell (average of lex2 with all neighbors).
func exchangeDecode(dst *data.Slice) {
	opencl.ExchangeDecode(dst, lex2.Gpu(), regions.Gpu(), M.Mesh())
}

// Set dst to the average dmi coupling per cell (average of din2 with all neighbors).
func dindDecode(dst *data.Slice) {
	opencl.ExchangeDecode(dst, din2.Gpu(), regions.Gpu(), M.Mesh())
}

// Returns the current exchange energy in Joules.
func GetExchangeEnergy() float64 {
	return -0.5 * cellVolume() * dot(&M_full, &B_exch)
}

// Scales the heisenberg exchange interaction between region1 and 2.
// Scale = 1 means the harmonic mean over the regions of Aex.
func ScaleInterExchange(region1, region2 int, scale float64) {
	lex2.setScale(region1, region2, scale)
}

// Sets the exchange interaction between region 1 and 2.
func InterExchange(region1, region2 int, value float64) {
	lex2.setInter(region1, region2, value)
}

// Scales the DMI interaction between region 1 and 2.
func ScaleInterDind(region1, region2 int, scale float64) {
	din2.setScale(region1, region2, scale)
}

// Sets the DMI interaction between region 1 and 2.
func InterDind(region1, region2 int, value float64) {
	din2.setInter(region1, region2, value)
}

// stores interregion exchange stiffness and DMI
// the interregion exchange/DMI by default is the harmonic mean (scale=1, inter=0)
type exchParam struct {
	parent         *RegionwiseScalar
	lut            [NREGION * (NREGION + 1) / 2]float64 // harmonic mean of regions (i,j)
	scale          [NREGION * (NREGION + 1) / 2]float64 // extra scale factor for lut[SymmIdx(i, j)]
	inter          [NREGION * (NREGION + 1) / 2]float64 // extra term for lut[SymmIdx(i, j)]
	gpu            opencl.SymmLUT                       // gpu copy of lut, lazily transferred when needed
	gpu_ok, cpu_ok bool                                 // gpu cache up-to date with lut source
}

// to be called after Aex or scaling changed
func (p *exchParam) invalidate() {
	p.cpu_ok = false
	p.gpu_ok = false
}

func (p *exchParam) init(parent *RegionwiseScalar) {
	for i := range p.scale {
		p.scale[i] = 1 // default scaling
		p.inter[i] = 0 // default additional interexchange term
	}
	p.parent = parent
}

// Get a GPU mirror of the look-up table.
// Copies to GPU first only if needed.
func (p *exchParam) Gpu() opencl.SymmLUT {
	p.update()
	if !p.gpu_ok {
		p.upload()
	}
	return p.gpu
}

// sets the interregion exchange/DMI using a specified value (scale = 0)
func (p *exchParam) setInter(region1, region2 int, value float64) {
	p.scale[symmidx(region1, region2)] = float64(0.)
	p.inter[symmidx(region1, region2)] = float64(value)
	p.invalidate()
}

// sets the interregion exchange/DMI by rescaling the harmonic mean (inter = 0)
func (p *exchParam) setScale(region1, region2 int, scale float64) {
	p.scale[symmidx(region1, region2)] = float64(scale)
	p.inter[symmidx(region1, region2)] = float64(0.)
	p.invalidate()
}

func (p *exchParam) update() {
	if !p.cpu_ok {
		ex := p.parent.cpuLUT()
		for i := 0; i < NREGION; i++ {
			exi := ex[0][i]
			for j := i; j < NREGION; j++ {
				exj := ex[0][j]
				I := symmidx(i, j)
				p.lut[I] = p.scale[I]*exchAverage(exi, exj) + p.inter[I]
			}
		}
		p.gpu_ok = false
		p.cpu_ok = true
	}
}

func (p *exchParam) upload() {
	// alloc if  needed
	if p.gpu == nil {
		p.gpu = opencl.SymmLUT(opencl.MemAlloc(len(p.lut) * opencl.SIZEOF_FLOAT64))
	}
	lut := p.lut // Copy, to work around Go 1.6 cgo pointer limitations.
	opencl.MemCpyHtoD(unsafe.Pointer(p.gpu), unsafe.Pointer(&lut[0]), opencl.SIZEOF_FLOAT64*len(p.lut))
	p.gpu_ok = true
}

// Index in symmetric matrix where only one half is stored.
// (!) Code duplicated in exchange.cu
func symmidx(i, j int) int {
	if j <= i {
		return i*(i+1)/2 + j
	} else {
		return j*(j+1)/2 + i
	}
}

// Returns the intermediate value of two exchange/dmi strengths.
// If both arguments have the same sign, the average mean is returned. If the arguments differ in sign
// (which is possible in the case of DMI), the geometric mean of the geometric and arithmetic mean is
// used. This average is continuous everywhere, monotonic increasing, and bounded by the argument values.
func exchAverage(exi, exj float64) float64 {
	if exi*exj >= 0.0 {
		return 2 / (1/exi + 1/exj)
	} else {
		exi_, exj_ := float64(exi), float64(exj)
		sign := math.Copysign(1, exi_+exj_)
		magn := math.Sqrt(math.Sqrt(-exi_*exj_) * math.Abs(exi_+exj_) / 2)
		return float64(sign * magn)
	}
}