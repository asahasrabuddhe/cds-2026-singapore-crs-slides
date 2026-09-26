for _, c := range keep.ToDrop() {
	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRCTL, prCapbsetDrop, uintptr(c), 0, 0, 0, 0)
	switch errno {
	case 0:
		dropped++
	case syscall.EINVAL:
		// This kernel does not know a capability with this number.
		// Nothing to drop, then.
	default:
		return dropped, fmt.Errorf("prctl(PR_CAPBSET_DROP, %s): %w", Name(c), errno)
	}
}
