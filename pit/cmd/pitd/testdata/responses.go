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
							"templated": true
						},
						{ 
							"id": "new", 
							"name": "links",
							"rel": [ "add" ], 
							"url": "/projects/", 
							"action": "append",
							"model": "n={name}\u0026d={description}"
						} 
					] 
				},
				{
					"id": "projects"
				}
			] 
		}
	}`
	OneProjectList = `
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
							"templated": true
						},
						{
							"id": "new",
							"name": "links",
							"rel": [ "add" ],
							"url": "/projects/",
							"action": "append",
							"model": "n={name}\u0026d={description}"
						}
					]
				},
				{
					"id": "projects",
					"data":
					[
						{
							"id": "101",
							"name": "project one",
							"rel": [ "self" ],
							"url": "/project/101",
							"data":
							[
								{ "rel": [ "add" ], "url": "/project/101/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026r={reporter}" },
								{ "rel": [ "search" ], "url": "/project/101/search{?name}", "templated": true}
							]
						}
					]
				}
			]
		}
	}`
)
