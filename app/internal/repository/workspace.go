package repository

import  (
	"gorm.io/gorm"
	
	"backend/internal/model"
)


type WorkspaceRepository interface {
	GetByUserID(userID uint) ([]model.Workspace, error)
	Create(workspace *model.Workspace) error
	Delete(id uint) error
	AddDocumentToWorkspace(documentID, workspaceID uint) error
	RemoveDocumentFromWorkspace(documentID uint) error
}

type workspaceRepo struct {
	db *gorm.DB
}


func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepo{db}
}

func (r *workspaceRepo) GetByUserID(userID uint) ([]model.Workspace, error) {
	var workspaces []model.Workspace
	err := r.db.Where("user_id = ?", userID).Find(&workspaces).Error
	return workspaces, err
}

func (r *workspaceRepo) Create(ws *model.Workspace) error {
	return r.db.Create(ws).Error
}

func (r *workspaceRepo) Delete(id uint) error {
	return r.db.Delete(&model.Workspace{}, id).Error
}

func (r *workspaceRepo) AddDocumentToWorkspace(documentID, workspaceID uint) error {
	return r.db.Model(&model.Document{}).Where("id = ?", documentID).Update("workspace_id", workspaceID).Error
}

func (r *workspaceRepo) RemoveDocumentFromWorkspace(documentID uint) error {
	return r.db.Model(&model.Document{}).Where("id = ?", documentID).Update("workspace_id", nil).Error
}