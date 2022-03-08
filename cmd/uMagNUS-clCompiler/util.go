package main

func generateCompilerOpts() string {
	return *Flag_ComArgs + " -std=" + *Flag_ClStd + " " + *Flag_includes + " " + *Flag_defines
}

func generateLinkerOprts() string {
	return ""
}
