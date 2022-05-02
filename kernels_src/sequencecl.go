package kernels_src

// Use the following lists to sequence order of file
// loads in order to build OpenCL Program
var OCLHeadersList = []string{
	"typedefs",
	"constants",
	"stdint",
	"stencil",
	"float3",
	"exchange",
	"atomicf",
	"reduce",
	"amul",
	"RNG_common",
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
	"reducesum_onestage",
	"reducesum_onestage_oop",
	"reducesum_onestage_inp",
	"reducesum_twophase",
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
	"xorwow_seed",
	"xorwow_uint",
	"xorwow_uniform",
	"xorwow_normal",
	"threefry_seed",
	"threefry_uint",
	"threefry_uniform",
	"threefry_normal"}

var OCL64KernelsList = []string{
	"xorwow64_seed",
	"xorwow64_ulong",
	"xorwow64_normal",
	"xorwow64_uniform",
	"threefry64_seed",
	"threefry64_ulong",
	"threefry64_normal",
	"threefry64_uniform"}
