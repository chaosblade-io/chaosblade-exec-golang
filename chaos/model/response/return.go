package response

func ReturnExperimentNotFound() Response {
	return Response{Code: ExperimentNotFound, Success: false, Error: "the experiment not found"}
}

func ReturnExperimentInCircuit() Response {
	return Response{Code: ExperimentInCircuit, Success: false, Error: "the current chaos state is in circuit"}
}

func ReturnExperimentMatcherNotFound(matcher string) Response {
	return Response{Code: ExperimentMatcherNotFound, Success: false, Error: "the " + matcher + " not found"}
}

func ReturnExperimentNotMatched(matcher string) Response {
	return Response{Code: ExperimentNotMatched, Success: false, Error: "the " + matcher + " not matched"}
}

func ReturnExperimentLimited() Response {
	return Response{Code: ExperimentLimited, Success: false, Error: "the experiment effect is limited"}
}

func ReturnIllegalParameters(error string) Response {
	return ReturnFail(IllegalParameters, error)
}