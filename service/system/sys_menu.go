package system

import (
	"errors"
	"github/shansec/go-vue-admin/dao/common/request"
	"github/shansec/go-vue-admin/global"
	"github/shansec/go-vue-admin/model/system"

	"gorm.io/gorm"
)

type MenuService struct{}

// CreateMenuService
// @author: [Shansec](https://github.com/shansec)
// @function: CreateMenuService
// @description: 添加菜单
// @param: menu system.SysBaseMenu
// @return: system.SysBaseMenu, error
func (menuService *MenuService) CreateMenuService(menu system.SysBaseMenu) (system.SysBaseMenu, error) {
	var menuInfo system.SysBaseMenu
	var err error
	if err = global.MAY_DB.Where("path = ? AND name = ?", menu.Path, menu.Name).First(&menuInfo).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return system.SysBaseMenu{}, errors.New("存在相同菜单")
	}
	err = global.MAY_DB.Create(&menu).Error
	if err != nil {
		return system.SysBaseMenu{}, err
	}
	return menu, nil
}

// DeleteMenuService
// @author: [Shansec](https://github.com/shansec)
// @function: DeleteMenuService
// @description: 删除菜单
// @param: menu system.SysBaseMenu
// @return: error
func (menuService *MenuService) DeleteMenuService(menu *system.SysBaseMenu) error {
	err := global.MAY_DB.Preload("SysRoles").Where("id = ?", menu.ID).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("菜单不存在！")
	}
	err = global.MAY_DB.Transaction(func(ctx *gorm.DB) error {
		var err error

		if err = ctx.Preload("SysRoles").Where("id = ?", menu.ID).First(menu).Error; err != nil {
			return err
		}

		if len(menu.SysRoles) != 0 {
			if err = ctx.Model(menu).Association("SysRoles").Delete(menu.SysRoles); err != nil {
				return err
			}
		}

		if err = ctx.Where("id = ?", menu.ID).Unscoped().Delete(menu).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

// GetMenuListService
// @author: [Shansec](https://github.com/shansec)
// @function: GetMenuListService
// @description: 分页获取菜单信息
// @param: pageInfo request.PageInfo
// @return: list interface{}, total int64, err error
func (menuService *MenuService) GetMenuListService(pageInfo request.PageInfo) (list interface{}, total int64, err error) {
	limit := pageInfo.PageSize
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	db := global.MAY_DB.Model(&system.SysBaseMenu{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var menuList []system.SysBaseMenu
	err = db.Limit(limit).Offset(offset).Where("parent_id = ?", "0").Find(&menuList).Error
	if err != nil {
		return nil, 0, err
	}
	for index := range menuList {
		err = menuService.findChildrenMenu(&menuList[index])
	}
	return menuList, total, nil
}

// findChildrenMenu
// @author: [Shansec](https://github.com/shansec)
// @function: findChildrenMenu
// @description: 分页获取菜单信息辅助方法，查找子菜单
// @param: menu *system.SysBaseMenu
// @return: err error
func (menuService *MenuService) findChildrenMenu(menu *system.SysBaseMenu) (err error) {
	err = global.MAY_DB.Where("parent_id = ?", menu.ID).Find(&menu.Children).Error
	if len(menu.Children) > 0 {
		for index := range menu.Children {
			err = menuService.findChildrenMenu(&menu.Children[index])
		}
	}
	return err
}
