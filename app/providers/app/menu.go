package app

type Menu struct {
	SubMenuId   int
	Title       string
	Link        string
	Icon        string
	SubMenu     []SubMenu
	HasSubmenu  bool
	IsSeparator bool
}

type SubMenu struct {
	Title string
	Link  string
}

var menu = []Menu{
	{
		Title:       "Dashboard",
		Link:        "/",
		Icon:        "icon-[tabler--home]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		SubMenuId:   1,
		Title:       "Scan",
		Link:        "/scan",
		Icon:        "icon-[tabler--scan]",
		HasSubmenu:  true,
		IsSeparator: false,
		SubMenu: []SubMenu{
			{
				Title: "E-Statement",
				Link:  "/scan/e-statement",
			},
			{
				Title: "ID Document",
				Link:  "/scan/id-document",
			},
			{
				Title: "Receipt & Invoice",
				Link:  "/scan/receipt-invoice",
			},
			{
				Title: "Business Card",
				Link:  "/scan/business-card",
			},
			{
				Title: "Barcode & QR-Code",
				Link:  "/scan/barcode-qrcode",
			},
		},
	},
	{
		Title:       "Billing",
		Link:        "/billing",
		Icon:        "icon-[tabler--receipt]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "Setting",
		Link:        "/setting",
		Icon:        "icon-[tabler--settings]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "Admin Zone",
		IsSeparator: true,
	},
	{
		SubMenuId:   2,
		Title:       "Master Data",
		Link:        "/master/",
		Icon:        "icon-[tabler--database-cog]",
		HasSubmenu:  true,
		IsSeparator: false,
		SubMenu: []SubMenu{
			{
				Title: "Theme",
				Link:  "/master/theme",
			},
			{
				Title: "Provider",
				Link:  "/master/provider",
			},
			{
				Title: "Payment Method",
				Link:  "/master/payment-method",
			},
		},
	},
	{
		SubMenuId:   3,
		Title:       "Manage User",
		Link:        "/manage-users/",
		Icon:        "icon-[tabler--users]",
		HasSubmenu:  true,
		IsSeparator: false,
		SubMenu: []SubMenu{
			{
				Title: "User",
				Link:  "/manage-users/user",
			},
			{
				Title: "Privilege",
				Link:  "/manage-users/privilege",
			},
		},
	},
	{
		Title:       "App Setting",
		Link:        "/app-setting",
		Icon:        "icon-[tabler--settings-code]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "Developer Zone",
		IsSeparator: true,
	},
	{
		Title:       "API Key",
		Link:        "/manage/api-key",
		Icon:        "icon-[tabler--key]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "Documentation",
		Link:        "/book/developer-documentation",
		Icon:        "icon-[tabler--files]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "API Documentation",
		Link:        "/book/api-documentation",
		Icon:        "icon-[tabler--files]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "Miscellaneous",
		IsSeparator: true,
	},
	{
		Title:       "Support",
		Link:        "/customer/support",
		Icon:        "icon-[tabler--headset]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
	{
		Title:       "User Guide",
		Link:        "/book/user-guide",
		Icon:        "icon-[tabler--files]",
		HasSubmenu:  false,
		IsSeparator: false,
	},
}
