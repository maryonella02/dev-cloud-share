package services

type ResourceService struct {
	BaseService
}

func NewResourceService(serviceURL string) *ResourceService {
	return &ResourceService{
		BaseService: BaseService{ServiceURL: serviceURL},
	}
}
