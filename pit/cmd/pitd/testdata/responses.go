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
							"id": 101,
							"name": "project one",
							"rel": [ "project" ],
							"data":
							[
								{ "rel": [ "issue" ], "url": "/projects/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026rr={reporter}" },
								{ "name": "id", "value": 101 },
								{ "name": "name", "value": "pone" },
								{ "name": "description", "value": "first test project" },
								{ "name": "owner", "value": "owner@test.net" }
							]
						}
					]
				}
			]
		}
	}`
)
