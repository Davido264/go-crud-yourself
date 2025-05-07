import { useState } from "react"
import { ArrowUpDown, Clock, Filter, Globe, HardDrive, MoreHorizontal, Search, Server, Settings } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

// Mock data for servers
const servers = [
  {
    id: 1,
    name: "Production API",
    ip: "192.168.1.1",
    location: "US East",
    status: "online",
    uptime: "99.9%",
    lastChecked: "2 minutes ago",
    type: "Production",
  },
  {
    id: 2,
    name: "Web Server",
    ip: "192.168.1.2",
    location: "US West",
    status: "online",
    uptime: "99.7%",
    lastChecked: "5 minutes ago",
    type: "Production",
  },
  {
    id: 3,
    name: "Database Cluster",
    ip: "192.168.1.3",
    location: "Europe",
    status: "offline",
    uptime: "98.2%",
    lastChecked: "1 minute ago",
    type: "Production",
  },
  {
    id: 4,
    name: "Staging API",
    ip: "192.168.1.4",
    location: "Asia",
    status: "online",
    uptime: "99.5%",
    lastChecked: "3 minutes ago",
    type: "Staging",
  },
  {
    id: 5,
    name: "Test Server",
    ip: "192.168.1.5",
    location: "Australia",
    status: "online",
    uptime: "99.1%",
    lastChecked: "10 minutes ago",
    type: "Development",
  },
  {
    id: 6,
    name: "Backup Server",
    ip: "192.168.1.6",
    location: "US East",
    status: "offline",
    uptime: "95.0%",
    lastChecked: "15 minutes ago",
    type: "Backup",
  },
]

export default function Dashboard() {
  const [searchQuery, setSearchQuery] = useState("")
  const [activeTab, setActiveTab] = useState("all")

  // Filter servers based on search query and active tab
  const filteredServers = servers.filter((server) => {
    const matchesSearch =
      server.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      server.ip.includes(searchQuery) ||
      server.location.toLowerCase().includes(searchQuery.toLowerCase())

    if (activeTab === "all") return matchesSearch
    if (activeTab === "online") return matchesSearch && server.status === "online"
    if (activeTab === "offline") return matchesSearch && server.status === "offline"
    if (activeTab === "production") return matchesSearch && server.type === "Production"

    return matchesSearch
  })

  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-10 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
        <div className="flex items-center gap-2">
          <Server className="h-6 w-6" />
          <h1 className="text-lg font-semibold">Server Monitor</h1>
        </div>
        <div className="ml-auto flex items-center gap-2">
          <Button variant="outline" size="icon">
            <Settings className="h-4 w-4" />
            <span className="sr-only">Settings</span>
          </Button>
        </div>
      </header>
      <main className="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-8">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle>Total Servers</CardTitle>
              <CardDescription>All registered servers</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{servers.length}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle>Online</CardTitle>
              <CardDescription>Servers currently online</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{servers.filter((s) => s.status === "online").length}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle>Offline</CardTitle>
              <CardDescription>Servers currently offline</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{servers.filter((s) => s.status === "offline").length}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle>Average Uptime</CardTitle>
              <CardDescription>Across all servers</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">98.7%</div>
            </CardContent>
          </Card>
        </div>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <h2 className="text-2xl font-bold tracking-tight">Servers</h2>
            <div className="flex items-center gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder="Search servers..."
                  className="w-full rounded-md pl-8 md:w-[300px]"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" size="sm" className="h-9 gap-1">
                    <Filter className="h-4 w-4" />
                    <span className="hidden sm:inline">Filter</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => setActiveTab("all")}>All Servers</DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setActiveTab("online")}>Online Only</DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setActiveTab("offline")}>Offline Only</DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setActiveTab("production")}>Production Only</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              <Button size="sm" className="h-9">
                <HardDrive className="mr-2 h-4 w-4" />
                Add Server
              </Button>
            </div>
          </div>
          <Tabs defaultValue="all" value={activeTab} onValueChange={setActiveTab}>
            <TabsList>
              <TabsTrigger value="all">All Servers</TabsTrigger>
              <TabsTrigger value="online">Online</TabsTrigger>
              <TabsTrigger value="offline">Offline</TabsTrigger>
              <TabsTrigger value="production">Production</TabsTrigger>
            </TabsList>
            <TabsContent value={activeTab} className="border rounded-md">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[100px]">Status</TableHead>
                    <TableHead>
                      <div className="flex items-center gap-1">
                        Name
                        <ArrowUpDown className="h-3 w-3" />
                      </div>
                    </TableHead>
                    <TableHead>IP Address</TableHead>
                    <TableHead>Location</TableHead>
                    <TableHead>Uptime</TableHead>
                    <TableHead>Last Check</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredServers.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-8 text-muted-foreground">
                        No servers found matching your criteria
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredServers.map((server) => (
                      <TableRow key={server.id}>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <div
                              className={`h-3 w-3 rounded-full ${server.status === "online" ? "bg-green-500" : "bg-red-500"}`}
                            />
                            <span className="capitalize">{server.status}</span>
                          </div>
                        </TableCell>
                        <TableCell className="font-medium">{server.name}</TableCell>
                        <TableCell>{server.ip}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Globe className="h-4 w-4 text-muted-foreground" />
                            {server.location}
                          </div>
                        </TableCell>
                        <TableCell>{server.uptime}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Clock className="h-4 w-4 text-muted-foreground" />
                            {server.lastChecked}
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant={server.type === "Production" ? "default" : "secondary"}>{server.type}</Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon">
                                <MoreHorizontal className="h-4 w-4" />
                                <span className="sr-only">Actions</span>
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem>View Details</DropdownMenuItem>
                              <DropdownMenuItem>Edit Server</DropdownMenuItem>
                              <DropdownMenuItem>Restart Server</DropdownMenuItem>
                              <DropdownMenuItem className="text-red-600">Remove Server</DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TabsContent>
          </Tabs>
        </div>
      </main>
    </div>
  )
}
