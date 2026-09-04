package workspace

// WorkspaceConfig tracks the configuration and state of a workspace
type WorkspaceConfig struct {
	RootPath       string
	GitInitialized bool
	SpecKitReady   bool
	DocsGenerated  []string
	CreatedFiles   []string
	GitRemoteURL   string
}

// NewWorkspaceConfig creates a new workspace configuration
func NewWorkspaceConfig(rootPath string) *WorkspaceConfig {
	return &WorkspaceConfig{
		RootPath:      rootPath,
		DocsGenerated: make([]string, 0),
		CreatedFiles:  make([]string, 0),
	}
}

// AddCreatedFile adds a file to the created files list
func (w *WorkspaceConfig) AddCreatedFile(path string) {
	w.CreatedFiles = append(w.CreatedFiles, path)
}

// AddGeneratedDoc adds a document to the generated docs list
func (w *WorkspaceConfig) AddGeneratedDoc(path string) {
	w.DocsGenerated = append(w.DocsGenerated, path)
}
