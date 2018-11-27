// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blacklist

var builtin = New(
	// Staff
	"abuse",
	"admin",
	"administration",
	"administrator",
	"administrators",
	"admins",
	"billing",
	"demo",
	"editor",
	"editors",
	"guest",
	"guests",
	"info",
	"marketplace",
	"master",
	"masters",
	"me",
	"member",
	"members",
	"moderator",
	"moderators",
	"myself",
	"owner",
	"owners",
	"partners",
	"root",
	"sales",
	"security",
	"self",
	"staff",
	"super",
	"superuser",
	"support",
	"sysadmin",
	"thethings",
	"things",
	"this",
	"tti",
	"ttn",
	"user",
	"users",
	"webmaster",

	// Reserved for subdomains
	"about",
	"account",
	"alerts",
	"alpha",
	"analytics",
	"app",
	"apps",
	"assets",
	"beta",
	"blog",
	"carreers",
	"cdn",
	"community",
	"compare",
	"conference",
	"contact",
	"customers",
	"dev",
	"dns",
	"docs",
	"events",
	"example",
	"examples",
	"extension",
	"extensions",
	"faq",
	"feedback",
	"forum",
	"git",
	"help",
	"home",
	"homes",
	"http",
	"https",
	"id",
	"identity",
	"join",
	"learn",
	"local",
	"localdomain",
	"localhost",
	"mail",
	"map",
	"metrics",
	"news",
	"notifications",
	"ops",
	"plugin",
	"plugins",
	"prod",
	"production",
	"profile",
	"secure",
	"ssh",
	"staging",
	"static",
	"statistics",
	"stats",
	"status",
	"team",
	"www",

	// TTN
	"applicationserver",
	"as",
	"broker",
	"cli",
	"console",
	"dashboard",
	"ga",
	"gatewayagent",
	"gatewayserver",
	"gs",
	"handler",
	"identityserver",
	"is",
	"joinserver",
	"js",
	"networkserver",
	"noc",
	"ns",
	"router",
	"ttnctl",
	"webui",
)
