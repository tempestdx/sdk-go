package resource

type CanonicalOperation struct {
	Type        CanonicalType
	Priority    int
	Concurrency int
}

type CanonicalType int

// TODO: cleanup and deduplicate canonical operations
const (
	Install CanonicalType = iota
	Uninstall
	Create
	Read
	Update
	Delete
	Upgrade
	Rollback
	Destroy
	Configure
	Write
	List
	Get
	Test
	Sync
)

func (r *ResourceDefinition) Install(fn OperationFunc) {
	r.RegisterOperation("_install", fn, Op.On(CanonicalOperation{Type: Install}))
}

func (r *ResourceDefinition) Create(fn OperationFunc) {
	r.RegisterOperation("_create", fn, Op.On(CanonicalOperation{Type: Create}))
}

// ... other methods
func ExecuteOperation() {

	// Route to appropriate operation

}
