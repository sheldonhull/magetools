package secrets

// // 🔐 Secrets scans for secret violations with gitleaks.
// func (Secrets) Check() error {
// 	checkPtermDebug()
// 	pterm.Info.Println("scans for secret violations")
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}
// 	c := []string{
// 		"run",
// 		"--rm",
// 		"-v",
// 		wd + ":/repo",
// 		"zricethezav/gitleaks:latest",
// 		"--commit=latest",
// 		"--no-git",
// 		"--format",
// 		"json",
// 		"--path=/repo",
// 		"--report=/repo/_artifacts/gitleaks.json",
// 		"--quiet",
// 	}
// 	if err := sh.Run("docker", c...); err != nil {
// 		return err
// 	}
// 	if _, err := os.Stat("_artifacts/gitleaks.json"); os.IsNotExist(err) {
// 		pterm.Success.Println("🔐  no _artifacts/gitleaks.json detected")
// 	} else {
// 		pterm.Error.Println("🔐  _artifacts/gitleaks.json detected")
// 		return errors.New("_artifacts/gitleaks.json detected")
// 	}

// 	return nil
// }
