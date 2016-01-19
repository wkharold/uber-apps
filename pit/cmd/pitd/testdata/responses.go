package testdata

var (
	EmptyTeamList = `
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
							"id": "project-list", 
							"name": "links",
							"rel": [ "collection" ], 
							"url": "/projects/", 
							"action": "read" 
						},
						{ 
							"id": "project-search", 
							"name": "links",
							"rel": [ "search" ], 
							"url": "/projects/search{?name}", 
							"templated": true,
							"action": "read"
						},
						{ 
							"id": "project-create", 
							"name": "links",
							"rel": [ "add" ], 
							"url": "/projects/", 
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
						}
					] 
				},
				{
					"id": "members"
				}
			] 
		}
	}`
	MultipleTeamMemberList = `
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
							"id": "project-list", 
							"name": "links",
							"rel": [ "collection" ], 
							"url": "/projects/", 
							"action": "read" 
						},
						{ 
							"id": "project-search", 
							"name": "links",
							"rel": [ "search" ], 
							"url": "/projects/search{?name}", 
							"templated": true,
							"action": "read"
						},
						{ 
							"id": "project-create", 
							"name": "links",
							"rel": [ "add" ], 
							"url": "/projects/", 
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
						}
					] 
				},
				{
					"id": "members",
					"data":
					[
						{
							"id": "1001",
							"rel": [ "self" ],
							"url": "/team/1001",
							"data":
							[
								{"name": "email", "value": "owner@test.net"}
							]
						}
					]
				},
				{
					"id": "members",
					"data":
					[
						{
							"id": "1002",
							"rel": [ "self" ],
							"url": "/team/1002",
							"data":
							[
								{"name": "email", "value": "owner@test.io"}
							]
						}
					]
				}
			] 
		}
	}`
	OneTeamMemberList = `
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
							"id": "project-list", 
							"name": "links",
							"rel": [ "collection" ], 
							"url": "/projects/", 
							"action": "read" 
						},
						{ 
							"id": "project-search", 
							"name": "links",
							"rel": [ "search" ], 
							"url": "/projects/search{?name}", 
							"templated": true,
							"action": "read"
						},
						{ 
							"id": "project-create", 
							"name": "links",
							"rel": [ "add" ], 
							"url": "/projects/", 
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
						}
					] 
				},
				{
					"id": "members",
					"data":
					[
						{
							"id": "1001",
							"rel": [ "self" ],
							"url": "/team/1001",
							"data":
							[
								{"name": "email", "value": "owner@test.net"}
							]
						}
					]
				}
			] 
		}
	}`
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
							"id": "project-list", 
							"name": "links",
							"rel": [ "collection" ], 
							"url": "/projects/", 
							"action": "read" 
						},
						{ 
							"id": "project-search", 
							"name": "links",
							"rel": [ "search" ], 
							"url": "/projects/search{?name}", 
							"templated": true,
							"action": "read"
						},
						{ 
							"id": "project-create", 
							"name": "links",
							"rel": [ "add" ], 
							"url": "/projects/", 
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
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
							"id": "project-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/projects/",
							"action": "read"
						},
						{
							"id": "project-search",
							"name": "links",
							"rel": [ "search" ],
							"url": "/projects/search{?name}",
							"templated": true,
							"action": "read"
						},
						{
							"id": "project-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/projects/",
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
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
							"id": "project-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/projects/",
							"action": "read"
						},
						{
							"id": "project-search",
							"name": "links",
							"rel": [ "search" ],
							"url": "/projects/search{?name}",
							"templated": true,
							"action": "read"
						},
						{
							"id": "project-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/projects/",
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
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
							"id": "project-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/projects/",
							"action": "read"
						},
						{
							"id": "project-search",
							"name": "links",
							"rel": [ "search" ],
							"url": "/projects/search{?name}",
							"templated": true,
							"action": "read"
						},
						{
							"id": "project-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/projects/",
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
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
							"id": "project-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/projects/",
							"action": "read"
						},
						{
							"id": "project-search",
							"name": "links",
							"rel": [ "search" ],
							"url": "/projects/search{?name}",
							"templated": true,
							"action": "read"
						},
						{
							"id": "project-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/projects/",
							"action": "append",
							"model": "n={name}\u0026d={description}\u0026o={owner}"
						},
						{
							"id": "team-members-list",
							"name": "links",
							"rel": [ "collection" ],
							"url": "/team",
							"action": "read"
						},
						{
							"id": "team-member-create",
							"name": "links",
							"rel": [ "add" ],
							"url": "/team",
							"action": "append",
							"model": "m={email}"
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
