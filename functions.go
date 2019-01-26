package main

func check_parameters(required int, parameters []string) bool {
	return len(parameters) == required
}

func help_function(_ []string) string{
	return help_text
}

func look_function(parameters []string) string{
	if check_parameters(1, parameters){

	} else {
		
	}
}