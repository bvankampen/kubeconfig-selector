package selector

func (s *Selector) GetHelpText() string {
	return "KS Help \n" +
		"Kubernetes Selector (ks) Version " + s.ctx.App.Version + "\n\n" +
		"  [yellow]q:[white]       Quit \n" +
		"  [yellow]<enter>:[white] Use Kubeconfig \n" +
		"  [yellow]m:[white]       Move Kubeconfig to " + s.appConfig.KubeconfigDir + " and use it " + "\n" +
		"  [yellow]d:[white]       Delete File \n" +
		"  [yellow]k:[white]       Toggle Kubeconfig \n" +
		"  [yellow](*):[white]     File not in " + s.appConfig.KubeconfigDir + "\n" +
		"  [yellow]?:[white]       Help \n" +
		" \n [green]Changelog: \n\n" +
		" [red]Version 1.1:[white] \n" +
		"   - Delete a kubeconfig file \n" +
		"   - This Help Screen \n"
}
