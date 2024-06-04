package limit

import (
	"hcm/cmd/data-service/service/capability"
	"hcm/pkg/cryptography"
	"hcm/pkg/dal/dao"
	"hcm/pkg/rest"
)

// InitService initial the limitRule service
func InitService(cap *capability.Capability) {
	svc := &service{
		dao:    cap.Dao,
		cipher: cap.Cipher,
	}

	h := rest.NewHandler()

	h.Add("ListLimitRule", "POST", "/limit_rule/list", svc.ListLimitRule)

	h.Load(cap.WebService)
}

type service struct {
	dao    dao.Set
	cipher cryptography.Crypto
}
