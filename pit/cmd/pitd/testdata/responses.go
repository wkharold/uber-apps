package testdata

var (
	EmptyProjectList = `
	{ 
		"uber": 
		{ 
			"version": "1.0", 
			"data": 
			[
				{ 
					"id": "links", 
					"data": 
					[ 
						{
							"id": "alps",
							"rel": [ "profile" ],
							"url": "/pit-alps.xml",
							"action": "read"
						},
						{ 
							"id": "list", 
							"name": "links",
							"rel": [ "collection" ], 
							"url": "/projects/", 
							"action": "read" 
						},
						{ 
							"id": "search", 
							"name": "links",
							"rel": [ "search" ], 
							"url": "/projects/search{?name}", 
							"templated": "true"
						},
						{ 
							"id": "new", 
							"name": "links",
							"rel": [ "add" ], 
							"url": "/projects/", 
							"action": "append",
							"model": "n={name}&d={description}"
						} 
					] 
				},
				{
					"id": "projects"
				}
			] 
		}
	}`
)
