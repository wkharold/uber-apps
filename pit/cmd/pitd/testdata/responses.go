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
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						} 
					] 
				},
				{
					"id": "projects"
				}
			] 
		}
	}`
	MultiProjectList = `
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
							"model": "n={name}\u0026d={description}\u0026o={owner}"
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
						},
						{
							"id": "102",
							"name": "project two",
							"rel": [ "self" ],
							"url": "/project/102",
							"data":
							[
								{ "rel": [ "add" ], "url": "/project/102/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026r={reporter}" },
								{ "rel": [ "search" ], "url": "/project/102/search{?name}", "templated": true}
							]
						},
						{
							"id": "103",
							"name": "project three",
							"rel": [ "self" ],
							"url": "/project/103",
							"data":
							[
								{ "rel": [ "add" ], "url": "/project/103/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026r={reporter}" },
								{ "rel": [ "search" ], "url": "/project/103/search{?name}", "templated": true}
							]
						}
					]
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
							"model": "n={name}\u0026d={description}\u0026o={owner}"
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
	Project101 = `
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
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						}
					]
				},
				{
					"id": "project",
					"data":
					[
						{
							"id": "101",
							"name": "project one",
							"rel": [ "self" ],
							"url": "/project/101",
							"data":
							[
								{"rel": [ "add" ], "url": "/project/101/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026r={reporter}"},
								{"rel": [ "search" ], "url": "/project/101/search{?name}", "templated": true},
								{"name": "description", "value": "first test project"},
								{"name": "owner", "value": "owner@test.net"}
							]
						}
					]
				}
			]
		}
	}`
	Project102 = `
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
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						}
					]
				},
				{
					"id": "project",
					"data":
					[
						{
							"id": "102",
							"name": "project two",
							"rel": [ "self" ],
							"url": "/project/102",
							"data":
							[
								{"rel": [ "add" ], "url": "/project/102/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026r={reporter}"},
								{"rel": [ "search" ], "url": "/project/102/search{?name}", "templated": true},
								{"name": "description", "value": "second test project"},
								{"name": "owner", "value": "owner@test.net"}
							]
						}
					]
				}
			]
		}
	}`
	UnknownProjectError = `
	{
		"uber":
		{
			"version": "1.0",
			"error":
			{
				"data":
				[
					{"name": "RequestFailed", "rel": ["reason"], "value": "No project exists with specified ID: [1]"}
				]
			}
		}
	}`
)
