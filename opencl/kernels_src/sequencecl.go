package kernels_src

// Use the following lists to sequence order of file
// loads in order to build OpenCL Program
var OCLHeadersList = []string{
	"constants",
	"stdint",
	"stencil",
	"float3",
	"exchange",
	"atomicf",
	"reduce",
	"amul",
	"RNG_common",
	"RNGmrg32k3a",
	"RNGmtgp",
	"RNGthreefry",
	"RNGxorwow",
	"sum"}

var OCLKernelsList = []string{
	"copypadmul2",
	"copyunpad",
	"crop",
	"addcubicanisotropy2",
	"pointwise_div",
	"divide",
	"adddmi",
	"adddmibulk",
	"dotproduct",
	"crossproduct",
	"addexchange",
	"exchangedecode",
	"kernmulC",
	"kernmulRSymm2Dxy",
	"kernmulRSymm2Dz",
	"kernmulRSymm3D",
	"llnoprecess",
	"lltorque2",
	"madd2",
	"madd3",
	"madd4",
	"madd5",
	"madd6",
	"madd7",
	"setmaxangle",
	"minimize",
	"mul",
	"normalize2",
	"addoommfslonczewskitorque",
	"addtworegionoommfslonczewskitorque",
	"reducedot",
	"reducemaxabs",
	"reducemaxdiff",
	"reducemaxvecdiff2",
	"reducemaxvecnorm2",
	"reducesum",
	"regionaddv",
	"regionadds",
	"regiondecode",
	"regionselect",
	"resize",
	"shiftbytes",
	"shiftbytesy",
	"shiftx",
	"shifty",
	"shiftz",
	"addmagnetoelasticfield",
	"getmagnetoelasticforce",
	"addslonczewskitorque2",
	"tworegionexchange_field",
	"tworegionexchange_edens",
	"setPhi",
	"setTheta",
	"settemperature2",
	"settopologicalcharge",
	"settopologicalchargelattice",
	"adduniaxialanisotropy",
	"adduniaxialanisotropy2",
	"addvoltagecontrolledanisotropy2",
	"vecnorm",
	"zeromask",
	"addzhanglitorque2",
	"mrg32k3a",
	"mtgp32_init_seed_kernel",
	"mtgp32_uniform",
	"mtgp32_normal",
	"xorwow_seed",
	"xorwow_uint",
	"xorwow_uniform",
	"xorwow_normal",
	"threefry_seed",
	"threefry_uint",
	"threefry_uniform",
	"threefry_normal",
	"square"}