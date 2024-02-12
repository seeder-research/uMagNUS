package engine

// Add arbitrary terms to B_eff, Edens_total.

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

var (
	B_custom       = NewVectorField("B_custom", "T", "User-defined field", AddCustomField)
	Edens_custom   = NewScalarField("Edens_custom", "J/m3", "Energy density of user-defined field.", AddCustomEnergyDensity)
	E_custom       = NewScalarValue("E_custom", "J", "total energy of user-defined field", GetCustomEnergy)
	customTerms    []Quantity // vector
	customEnergies []Quantity // scalar
)

func init() {
	registerEnergy(GetCustomEnergy, AddCustomEnergyDensity)
	DeclFunc("AddFieldTerm", AddFieldTerm, "Add an expression to B_eff.")
	DeclFunc("AddEdensTerm", AddEdensTerm, "Add an expression to Edens.")
	DeclFunc("Add", Add, "Add two quantities")
	DeclFunc("Madd", Madd, "Weighted addition: Madd(Q1,Q2,c1,c2) = c1*Q1 + c2*Q2")
	DeclFunc("Dot", Dot, "Dot product of two vector quantities")
	DeclFunc("Cross", Cross, "Cross product of two vector quantities")
	DeclFunc("Mul", Mul, "Point-wise product of two quantities")
	DeclFunc("MulMV", MulMV, "Matrix-Vector product: MulMV(AX, AY, AZ, m) = (AX·m, AY·m, AZ·m). "+
		"The arguments Ax, Ay, Az and m are quantities with 3 componets.")
	DeclFunc("Div", Div, "Point-wise division of two quantities")
	DeclFunc("Const", Const, "Constant, uniform number")
	DeclFunc("ConstVector", ConstVector, "Constant, uniform vector")
	DeclFunc("Shifted", Shifted, "Shifted quantity")
	DeclFunc("Masked", Masked, "Mask quantity with shape")
	DeclFunc("Normalized", Normalized, "Normalize quantity")
	DeclFunc("CustomQuantity", CustomQuantity, "Custom scalar/vector quantity defined using array")
	DeclFunc("RemoveCustomFields", RemoveCustomFields, "Removes all custom fields again")
	DeclFunc("RemoveCustomEnergies", RemoveCustomEnergies, "Removes all custom energies")
}

// Removes all customfields
func RemoveCustomFields() {
	customTerms = nil
}

// Removes all customenergies
func RemoveCustomEnergies() {
	customEnergies = nil
}

// AddFieldTerm adds an effective field function (returning Teslas) to B_eff.
// Be sure to also add the corresponding energy term using AddEnergyTerm.
func AddFieldTerm(b Quantity) {
	customTerms = append(customTerms, b)
}

// AddEnergyTerm adds an energy density function (returning Joules/m³) to Edens_total.
// Needed when AddFieldTerm was used and a correct energy is needed
// (e.g. for Relax, Minimize, ...).
func AddEdensTerm(e Quantity) {
	customEnergies = append(customEnergies, e)
}

// AddCustomField evaluates the user-defined custom field terms
// and adds the result to dst.
func AddCustomField(dst *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in addcustomfield: %+v \n", err)
	}
	q1Idx, q2Idx, q3Idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1Idx)
	defer opencl.CheckinQueue(q2Idx)
	defer opencl.CheckinQueue(q3Idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1Idx], opencl.ClCmdQueue[q2Idx], opencl.ClCmdQueue[q3Idx]}
	for _, term := range customTerms {
		buf := ValueOf(term)
		opencl.Add(dst, dst, buf, queues, nil)
		// sync every iteration due to recycle
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("error waiting for all queues to finish in addcustomfield: %+v \n", err)
		}
		opencl.Recycle(buf)
	}
}

// Adds the custom energy densities (defined with AddCustomE
func AddCustomEnergyDensity(dst *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in addcustomfield: %+v \n", err)
	}
	seqQueue := opencl.ClCmdQueue[0]
	for _, term := range customEnergies {
		buf := ValueOf(term)
		opencl.Add(dst, dst, buf, []*cl.CommandQueue{seqQueue}, nil)
		// sync every iteration due to recycle
		if err := opencl.WaitAllQueuesToFinish(); err != nil {
			fmt.Printf("error waiting for all queues to finish in addcustomfield: %+v \n", err)
		}
		opencl.Recycle(buf)
	}
}

func GetCustomEnergy() float64 {
	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	AddCustomEnergyDensity(buf)
	seqQueue := opencl.ClCmdQueue[0]
	return cellVolume() * float64(opencl.Sum(buf, seqQueue, nil))
}

type constValue struct {
	value []float64
}

func (c *constValue) NComp() int { return len(c.value) }

func (d *constValue) EvalTo(dst *data.Slice) {
	for c, v := range d.value {
		opencl.Memset(dst.Comp(c), float32(v))
	}
}

// Const returns a constant (uniform) scalar quantity,
// that can be used to construct custom field terms.
func Const(v float64) Quantity {
	return &constValue{[]float64{v}}
}

// ConstVector returns a constant (uniform) vector quantity,
// that can be used to construct custom field terms.
func ConstVector(x, y, z float64) Quantity {
	return &constValue{[]float64{x, y, z}}
}

// fieldOp holds the abstract functionality for operations
// (like add, multiply, ...) on space-dependend quantites
// (like M, B_sat, ...)
type fieldOp struct {
	a, b  Quantity
	nComp int
}

func (o fieldOp) NComp() int {
	return o.nComp
}

type dotProduct struct {
	fieldOp
}

type crossProduct struct {
	fieldOp
}

type addition struct {
	fieldOp
}

type mAddition struct {
	fieldOp
	fac1, fac2 float64
}

type mulmv struct {
	ax, ay, az, b Quantity
}

// MulMV returns a new Quantity that evaluates to the
// matrix-vector product (Ax·b, Ay·b, Az·b).
func MulMV(Ax, Ay, Az, b Quantity) Quantity {
	util.Argument(Ax.NComp() == 3 &&
		Ay.NComp() == 3 &&
		Az.NComp() == 3 &&
		b.NComp() == 3)
	return &mulmv{Ax, Ay, Az, b}
}

func (q *mulmv) EvalTo(dst *data.Slice) {
	util.Argument(dst.NComp() == 3)
	opencl.Zero(dst)
	b := ValueOf(q.b)
	defer opencl.Recycle(b)

	// sync in the beginning
	var err error
	if err = opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in mulmv.evalto(): %+v \n", err)
	}
	seqQueue := opencl.ClCmdQueue[0]
	{
		Ax := ValueOf(q.ax)
		opencl.AddDotProduct(dst.Comp(X), 1, Ax, b, seqQueue, nil)
		if err = seqQueue.Finish(); err != nil {
			fmt.Printf("error waiting for ax to compute: %+v \n", err)
		}
		opencl.Recycle(Ax)
	}
	{

		Ay := ValueOf(q.ay)
		opencl.AddDotProduct(dst.Comp(Y), 1, Ay, b, seqQueue, nil)
		if err = seqQueue.Finish(); err != nil {
			fmt.Printf("error waiting for ay to compute: %+v \n", err)
		}
		opencl.Recycle(Ay)
	}
	{
		Az := ValueOf(q.az)
		opencl.AddDotProduct(dst.Comp(Z), 1, Az, b, seqQueue, nil)
		if err = seqQueue.Finish(); err != nil {
			fmt.Printf("error waiting for az to compute: %+v \n", err)
		}
		opencl.Recycle(Az)
	}
}

func (q *mulmv) NComp() int {
	return 3
}

// DotProduct creates a new quantity that is the dot product of
// quantities a and b. E.g.:
//
//	DotProct(&M, &B_ext)
func Dot(a, b Quantity) Quantity {
	return &dotProduct{fieldOp{a, b, 1}}
}

func (d *dotProduct) EvalTo(dst *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in dotproduct.evalto(): %+v \n", err)
	}
	A := ValueOf(d.a)
	defer opencl.Recycle(A)
	B := ValueOf(d.b)
	defer opencl.Recycle(B)
	opencl.Zero(dst)
	seqQueue := opencl.ClCmdQueue[0]
	opencl.AddDotProduct(dst, 1, A, B, seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for seqQueue to finish in dotproduct.evalto(): %+v \n", err)
	}
}

// CrossProduct creates a new quantity that is the cross product of
// quantities a and b. E.g.:
//
//	CrossProct(&M, &B_ext)
func Cross(a, b Quantity) Quantity {
	return &crossProduct{fieldOp{a, b, 3}}
}

func (d *crossProduct) EvalTo(dst *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in dotproduct.evalto(): %+v \n", err)
	}
	A := ValueOf(d.a)
	defer opencl.Recycle(A)
	B := ValueOf(d.b)
	defer opencl.Recycle(B)
	opencl.Zero(dst)
	seqQueue := opencl.ClCmdQueue[0]
	opencl.CrossProduct(dst, A, B, seqQueue, nil)
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for seqQueue to finish in dotproduct.evalto(): %+v \n", err)
	}
}

func Add(a, b Quantity) Quantity {
	if a.NComp() != b.NComp() {
		panic(fmt.Sprintf("Cannot point-wise Add %v components by %v components", a.NComp(), b.NComp()))
	}
	return &addition{fieldOp{a, b, a.NComp()}}
}

func (d *addition) EvalTo(dst *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before dotproduct.evalto(): %+v \n", err)
	}
	A := ValueOf(d.a)
	defer opencl.Recycle(A)
	B := ValueOf(d.b)
	defer opencl.Recycle(B)
	opencl.Zero(dst)
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	seqQueue := opencl.ClCmdQueue[0]
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
	opencl.Add(dst, A, B, queues, nil)
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after dotproduct.evalto(): %+v \n", err)
	}
}

type pointwiseMul struct {
	fieldOp
}

func Madd(a, b Quantity, fac1, fac2 float64) *mAddition {
	if a.NComp() != b.NComp() {
		panic(fmt.Sprintf("Cannot point-wise add %v components by %v components", a.NComp(), b.NComp()))
	}
	return &mAddition{fieldOp{a, b, a.NComp()}, fac1, fac2}
}

func (o *mAddition) EvalTo(dst *data.Slice) {
	A := ValueOf(o.a)
	defer opencl.Recycle(A)
	B := ValueOf(o.b)
	defer opencl.Recycle(B)
	opencl.Zero(dst)
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	seqQueue := opencl.ClCmdQueue[0]
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
	opencl.Madd2(dst, A, B, float32(o.fac1), float32(o.fac2), queues, nil)
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after dotproduct.evalto(): %+v \n", err)
	}
}

// Mul returns a new quantity that evaluates to the pointwise product a and b.
func Mul(a, b Quantity) Quantity {
	nComp := -1
	switch {
	case a.NComp() == b.NComp():
		nComp = a.NComp() // vector*vector, scalar*scalar
	case a.NComp() == 1:
		nComp = b.NComp() // scalar*something
	case b.NComp() == 1:
		nComp = a.NComp() // something*scalar
	default:
		panic(fmt.Sprintf("Cannot point-wise multiply %v components by %v components", a.NComp(), b.NComp()))
	}

	return &pointwiseMul{fieldOp{a, b, nComp}}
}

func (d *pointwiseMul) EvalTo(dst *data.Slice) {
	opencl.Zero(dst)
	a := ValueOf(d.a)
	defer opencl.Recycle(a)
	b := ValueOf(d.b)
	defer opencl.Recycle(b)

	switch {
	case a.NComp() == b.NComp():
		mulNN(dst, a, b) // vector*vector, scalar*scalar
	case a.NComp() == 1:
		mul1N(dst, a, b)
	case b.NComp() == 1:
		mul1N(dst, b, a)
	default:
		panic(fmt.Sprintf("Cannot point-wise multiply %v components by %v components", a.NComp(), b.NComp()))
	}
}

// mulNN pointwise multiplies two N-component vectors,
// yielding an N-component vector stored in dst.
func mulNN(dst, a, b *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before mulnn: %+v \n", err)
	}
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	seqQueue := opencl.ClCmdQueue[0]
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
	opencl.Mul(dst, a, b, queues, nil)
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after mulnn: %+v \n", err)
	}
}

// mul1N pointwise multiplies a scalar (1-component) with an N-component vector,
// yielding an N-component vector stored in dst.
func mul1N(dst, a, b *data.Slice) {
	util.Assert(a.NComp() == 1)
	util.Assert(dst.NComp() == b.NComp())
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before mul1n: %+v \n", err)
	}
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	seqQueue := opencl.ClCmdQueue[0]
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
	for c := 0; c < dst.NComp(); c++ {
		opencl.Mul(dst.Comp(c), a, b.Comp(c), []*cl.CommandQueue{queues[c]}, nil)
	}
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after mul1n: %+v \n", err)
	}
}

type pointwiseDiv struct {
	fieldOp
}

// Div returns a new quantity that evaluates to the pointwise product a and b.
func Div(a, b Quantity) Quantity {
	nComp := -1
	switch {
	case a.NComp() == b.NComp():
		nComp = a.NComp() // vector/vector, scalar/scalar
	case b.NComp() == 1:
		nComp = a.NComp() // something/scalar
	default:
		panic(fmt.Sprintf("Cannot point-wise divide %v components by %v components", a.NComp(), b.NComp()))
	}
	return &pointwiseDiv{fieldOp{a, b, nComp}}
}

func (d *pointwiseDiv) EvalTo(dst *data.Slice) {
	a := ValueOf(d.a)
	defer opencl.Recycle(a)
	b := ValueOf(d.b)
	defer opencl.Recycle(b)

	switch {
	case a.NComp() == b.NComp():
		divNN(dst, a, b) // vector*vector, scalar*scalar
	case b.NComp() == 1:
		divN1(dst, a, b)
	default:
		panic(fmt.Sprintf("Cannot point-wise divide %v components by %v components", a.NComp(), b.NComp()))
	}

}

func divNN(dst, a, b *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before divnn: %+v \n", err)
	}
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	seqQueue := opencl.ClCmdQueue[0]
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
	opencl.Div(dst, a, b, queues, nil)
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after divnn: %+v \n", err)
	}
}

func divN1(dst, a, b *data.Slice) {
	util.Assert(dst.NComp() == a.NComp())
	util.Assert(b.NComp() == 1)
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before divn1: %+v \n", err)
	}
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	seqQueue := opencl.ClCmdQueue[0]
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
	for c := 0; c < dst.NComp(); c++ {
		opencl.Div(dst.Comp(c), a.Comp(c), b, []*cl.CommandQueue{queues[c]}, nil)
	}
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after divn1: %+v \n", err)
	}
}

type shifted struct {
	orig       Quantity
	dx, dy, dz int
}

// Shifted returns a new Quantity that evaluates to
// the original, shifted over dx, dy, dz cells.
func Shifted(q Quantity, dx, dy, dz int) Quantity {
	util.Assert(dx != 0 || dy != 0 || dz != 0)
	return &shifted{q, dx, dy, dz}
}

func (q *shifted) EvalTo(dst *data.Slice) {
	orig := ValueOf(q.orig)
	defer opencl.Recycle(orig)
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before shifted.evalto(): %+v \n", err)
	}
	queueIndices := make([]int, dst.NComp())
	queues := make([]*cl.CommandQueue, dst.NComp())
	for idx := 0; idx < len(queueIndices); idx++ {
		queueIndex := opencl.CheckoutQueue()
		queueIndices[idx] = queueIndex
		defer opencl.CheckinQueue(queueIndices[idx])
		queues[idx] = opencl.ClCmdQueue[queueIndex]
	}
	for i := 0; i < q.NComp(); i++ {
		dsti := dst.Comp(i)
		origi := orig.Comp(i)
		if q.dx != 0 {
			opencl.ShiftX(dsti, origi, q.dx, 0, 0, queues[i], nil)
		}
		if q.dy != 0 {
			opencl.ShiftY(dsti, origi, q.dy, 0, 0, queues[i], nil)
		}
		if q.dz != 0 {
			opencl.ShiftZ(dsti, origi, q.dz, 0, 0, queues[i], nil)
		}
	}
	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish after shifted.evalto(): %+v \n", err)
	}
}

func (q *shifted) NComp() int {
	return q.orig.NComp()
}

// Masks a quantity with a shape
// The shape will be only evaluated once on the mesh,
// and will be re-evaluated after mesh change,
// because otherwise too slow
func Masked(q Quantity, shape Shape) Quantity {
	return &masked{q, shape, nil, data.Mesh{}}
}

type masked struct {
	orig  Quantity
	shape Shape
	mask  *data.Slice
	mesh  data.Mesh
}

func (q *masked) EvalTo(dst *data.Slice) {
	if q.mesh != *Mesh() {
		// When mesh is changed, mask needs an update
		q.createMask()
	}
	orig := ValueOf(q.orig)
	defer opencl.Recycle(orig)
	mul1N(dst, q.mask, orig)
}

func (q *masked) NComp() int {
	return q.orig.NComp()
}

func (q *masked) createMask() {
	size := Mesh().Size()
	// Prepare mask on host
	maskhost := data.NewSlice(SCALAR, size)
	defer maskhost.Free()
	maskScalars := maskhost.Scalars()
	for iz := 0; iz < size[Z]; iz++ {
		for iy := 0; iy < size[Y]; iy++ {
			for ix := 0; ix < size[X]; ix++ {
				r := Index2Coord(ix, iy, iz)
				if q.shape(r[X], r[Y], r[Z]) {
					maskScalars[iz][iy][ix] = 1
				}
			}
		}
	}
	// Update mask
	q.mask.Free()
	q.mask = opencl.NewSlice(SCALAR, size)
	data.Copy(q.mask, maskhost)
	// sync copy because host data lives in function
	if err := opencl.H2DQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish in createmask: %+v \n", err)
	}
	q.mesh = *Mesh()
	// Remove mask from host
}

// Normalized returns a quantity that evaluates to the unit vector of q
func Normalized(q Quantity) Quantity {
	return &normalized{q}
}

type normalized struct {
	orig Quantity
}

func (q *normalized) NComp() int {
	return 3
}

func (q *normalized) EvalTo(dst *data.Slice) {
	util.Assert(dst.NComp() == q.NComp())
	q.orig.EvalTo(dst)
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish before normalize: %+v \n", err)
	}
	seqQueue := opencl.ClCmdQueue[0]
	opencl.Normalize(dst, nil, seqQueue, nil)
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for normalize to finish: %+v \n", err)
	}
}

func CustomQuantity(inSlice *data.Slice) Quantity {
	util.Assert(inSlice.NComp() == 1 || inSlice.NComp() == 3)
	size := Mesh().Size()
	sliceSize := inSlice.Size()
	util.Assert(size[X] == sliceSize[X] && size[Y] == sliceSize[Y] && size[Z] == sliceSize[Z])
	retQuant := &customQuantity{nil, size}
	// sync in the beginning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		fmt.Printf("error waiting for all queues to finish initializing customquantity 0: %+v \n", err0)
		fmt.Printf("error waiting for all queues to finish initializing customquantity 1: %+v \n", err1)
		fmt.Printf("error waiting for all queues to finish initializing customquantity 2: %+v \n", err2)
	}
	if inSlice.NComp() == 1 {
		retQuant.customquant = opencl.NewSlice(SCALAR, size)
	} else {
		retQuant.customquant = opencl.NewSlice(VECTOR, size)
	}
	data.Copy(retQuant.customquant, inSlice)
	// sync before returning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		fmt.Printf("error waiting for all queues to finish after customquantity 0: %+v \n", err0)
		fmt.Printf("error waiting for all queues to finish after customquantity 1: %+v \n", err1)
		fmt.Printf("error waiting for all queues to finish after customquantity 2: %+v \n", err2)
	}
	return retQuant
}

type customQuantity struct {
	customquant *data.Slice
	size        [3]int
}

func (q *customQuantity) NComp() int {
	return q.customquant.NComp()
}

func (q *customQuantity) EvalTo(dst *data.Slice) {
	util.Assert(dst.NComp() == q.customquant.NComp())
	// sync in the beginning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		fmt.Printf("error waiting for all queues to finish in customquantity.evalto() 0: %+v \n", err0)
		fmt.Printf("error waiting for all queues to finish in customquantity.evalto() 1: %+v \n", err1)
		fmt.Printf("error waiting for all queues to finish in customquantity.evalto() 2: %+v \n", err2)
	}
	data.Copy(dst, q.customquant)
	// sync before returning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		fmt.Printf("error waiting for all queues to finish after customquantity.evalto() 0: %+v \n", err0)
		fmt.Printf("error waiting for all queues to finish after customquantity.evalto() 1: %+v \n", err1)
		fmt.Printf("error waiting for all queues to finish after customquantity.evalto() 2: %+v \n", err2)
	}
}
