package main

func printCommands(cmd Command, level int) {
	for i := 0; i < level; i++ {
		print("  ")
	}
	println("Command:" + namePrint(cmd))
	for i := 0; i < level; i++ {
		print("  ")
	}
	println("Args:")
	for _, arg := range cmd.Args {
		for i := 0; i < level; i++ {
			print("  ")
		}
		println("  ", arg)
	}
	for i := 0; i < level; i++ {
		print("  ")
	}
	println("Flags:")
	for flag, value := range cmd.Flags {
		for i := 0; i < level; i++ {
			print("  ")
		}
		println("  ", flag, value)
	}
	for _, sub := range cmd.Sub {
		printCommands(*sub, level+1)
	}
}

func namePrint(cmd Command) string {
	// recursive print
	if cmd.Parent != nil {
		return namePrint(*cmd.Parent) + " " + cmd.Name
	}
	return cmd.Name
}
