# ubiwhereChallenge

Based on the previous requirements, develop a Go application that
implements the following features:
	- Collect CPU, RAM each second from the operating system and
store it in a local database;
	- Create a simulator for the external device, which must
generate random data samples. Each sample must be
composed by 4 variables (for instance, 4 integers).
	- Collect samples (each second) generated by the previously
mentioned simulator and store them in local database.
	- Provide an interface through the console allowing the
following commands:
		○ Get last n metrics for all variables
		○ Get last n metrics for one or more variables
		○ Get an average of the value of one or more variables
    
    
 To Start the application:
  ./ubiwhere (without arguments)
  ./ubiwhere 15 (with arguments)
 
 The argument represent the amount of seconds the values will be displayed.
