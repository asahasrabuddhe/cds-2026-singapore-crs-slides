// The same binary, again, as "mycontainer child". Everything that
// makes it a container and not a process happens inside clone(2),
// through SysProcAttr.
cmd := exec.Command("/proc/self/exe", ChildCommand)
cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
cmd.ExtraFiles = []*os.File{specRead} // the spec, as fd 3
cmd.SysProcAttr = sysProcAttr(s, group)

if err := cmd.Start(); err != nil {
	return 0, fmt.Errorf("start container: %w", err)
}
