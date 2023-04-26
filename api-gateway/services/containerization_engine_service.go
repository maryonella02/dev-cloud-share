package services

type ContainerizationEngineService struct {
	BaseService
}

func NewContainerizationEngineService(serviceURL string) *ContainerizationEngineService {
	return &ContainerizationEngineService{
		BaseService: BaseService{ServiceURL: serviceURL},
	}
}
