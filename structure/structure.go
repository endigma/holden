package structure

// Page is a struct that holds webpage info for the template
type Page struct {
	Prefix           string
	Contents         string
	Meta             map[string]interface{}
	SidebarContents  string
	Raw              string
	DisplayBackToTop bool
	DisplaySidebar   bool
}

// Directory is a struct that holds information
// about the directory that holds the served markdown files
type Directory struct {
	Directories []Directory
	Files       []string
	Name        string
}
