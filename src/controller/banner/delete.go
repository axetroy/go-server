package banner

import "github.com/axetroy/go-server/src/service"

func DeleteBannerById(id string) {
	service.DeleteRowByTable("banner", "id", id)
}
