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
						},
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
	IssuesProject102 = `
	{
		"id": "issues",
		"rel": [ "self" ],
		"url": "/project/102/issues",
		"action": "read",
		"data": []
	}`
	IssuesProject103 = `
	{
		"id": "issues",
		"rel": [ "self" ],
		"url": "/project/103/issues",
		"action": "read",
		"data": [
			{
				"id": "1031",
				"name": "issueone",
				"rel": [ "self" ],
				"url": "/project/103/issue/1031",
				"action": "read",
				"data":
				[
					{"rel": [ "close" ], "url": "/project/103/issue/close", "action": "append", "model": "i=1031"},
					{"rel": [ "return" ], "url": "/project/103/issue/return", "action": "append", "model": "i=1031"},
					{"rel": [ "assign" ], "url": "/project/103/issue/1031/assignements", "action": "append", "model": "m={member}"},d
					{"name": "description", "value": "issue one"},
					{"name": "priority", "value": "1"},
					{"name": "status", "value": "OPEN"},
					{"name": "reporter", "value": "fred@testrock.org"},
					{
						"id": "assignments",
						"rel": [ "self" ],
						"url": "/project/102/issue/1031/assignments",
						"action": "read",
						"data": []
					}
				]
			},
			{
				"id": "1032",
				"name": "issuetwo",
				"rel": [ "self" ],
				"url": "/project/103/issue/1032",
				"action": "read",
				"data":
				[
					{"rel": [ "close" ], "url": "/project/103/issue/close", "action": "append", "model": "i=1032"},
					{"rel": [ "return" ], "url": "/project/103/issue/return", "action": "append", "model": "i=1032"},
					{"rel": [ "assign" ], "url": "/project/103/issue/1032/assignements", "action": "append", "model": "m={member}"},d
					{"name": "description", "value": "issue two"},
					{"name": "priority", "value": "1"},
					{"name": "status", "value": "OPEN"},
					{"name": "reporter", "value": "fred@testrock.org"},
					{
						"id": "assignments",
						"rel": [ "self" ],
						"url": "/project/103/issue/1032/assignments",
						"action": "read",
						"data": []
					}
				]
			},
			{
				"id": "1033",
				"name": "issuethree",
				"rel": [ "self" ],
				"url": "/project/103/issue/1033",
				"action": "read",
				"data":
				[
					{"rel": [ "close" ], "url": "/project/103/issue/close", "action": "append", "model": "i=1033"},
					{"rel": [ "return" ], "url": "/project/103/issue/return", "action": "append", "model": "i=1033"},
					{"rel": [ "assign" ], "url": "/project/103/issue/1033/assignements", "action": "append", "model": "m={member}"},d
					{"name": "description", "value": "issue three"},
					{"name": "priority", "value": "3"},
					{"name": "status", "value": "CLOSED"},
					{"name": "reporter", "value": "barney@testrock.org"},
					{
						"id": "assignments",
						"rel": [ "self" ],
						"url": "/project/103/issue/1033/assignments",
						"action": "read",
						"data": []
					}
				]
			}
		]
	}`
	ProjectWithIssuesAndMembers = `
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
							"action": "read",
							"data":
							[
								{"rel": [ "add" ], "url": "/project/102/issues", "action": "append", "model": "n={name}\u0026d={description}\u0026p={priority}\u0026r={reporter}"},
								{"rel": [ "collection" ], "url": "/project/102/issues", "action": "read"},
								{"rel": [ "search" ], "url": "/project/102/search{?name}", "templated": true, "action": "read"},
								{"name": "description", "value": "second test project"},
								{"name": "owner", "value": "owner@test.net"}
								"data":
								[
									{
										"id": "issues",
										"rel": [ "self" ],
										"url": "/project/102/issues",
										"action": "read",
										"data":
										[
											{
												"id": "2001",
												"name": "issueone",
												"rel": [ "self" ],
												"url": "/project/102/issue/2001",
												"action": "read",
												"data":
												[
													{"rel": [ "close" ], "url": "/project/102/issue/close", "action": "append", "model": "i=2001"},
													{"rel": [ "return" ], "url": "/project/102/issue/return", "action": "append", "model": "i=2001"},
													{"rel": [ "assign" ], "url": "/project/102/issue/2001/assignements", "action": "append", "model": "m={member}"},d
													{"name": "description", "value": "issue one"},
													{"name": "priority", "value": "1"},
													{"name": "status", "value": "OPEN"},
													{"name": "reporter", "value": "fred@testrock.org"},
													{
														"id": "assignments",
														"rel": [ "self" ],
														"url": "/project/102/issue/2001/assignments",
														"action": "read",
														"data":
														[
															{
																{"rel": [ "remove" ], "url": "/project/102/issue/2001/assignment/1006", "action": "delete"},
																{"rel": [ "member" ], "url": "/team/1006", "action": "read"},
																{"name": "email", "value": "alice@members.com"}
															}
														]
													}
												]
											},
										]
									},
									{
										"id": "contributors",
										"rel": [ "self" ],
										"url": "/project/102/contributors",
										"action": "read",
										"data":
										[
											{
												"id": "1006",
												"rel": [ "self" ],
												"url": "/project/102/contributor/1006",
												"data":
												[
													{"rel": [ "unassign" ], "url": "/project/102/contributor/unassign", "action": "append", "model": "m=alice@members.com"},
													{"name": "email", "value": "alice@members.com"}
												]
											}
										]
									}
								]
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
