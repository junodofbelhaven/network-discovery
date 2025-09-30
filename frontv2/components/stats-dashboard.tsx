"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, PieChart, Pie, Cell, ResponsiveContainer } from "recharts"
import { Network, Wifi, Clock, Shield, Database } from "lucide-react"

interface StatsDashboardProps {
  statistics: {
    total_devices: number
    reachable_devices: number
    snmp_devices: number
    arp_only_devices: number
    devices_with_mac: number
    vendor_distribution: Record<string, number>
    scan_method_distribution: Record<string, number>
    avg_response_time_ms: number
  }
  scanInfo: {
    scan_type: string
    network_range: string
    snmp_communities: string[]
    timeout: number
    retries: number
    worker_count: number
  }
}

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#8884D8", "#82CA9D"]

export function StatsDashboard({ statistics, scanInfo }: StatsDashboardProps) {
  const vendorData = Object.entries(statistics.vendor_distribution).map(([name, value]) => ({
    name,
    value,
  }))

  const methodData = Object.entries(statistics.scan_method_distribution).map(([name, value]) => ({
    name,
    value,
  }))

  const reachabilityPercentage = (statistics.reachable_devices / statistics.total_devices) * 100
  const snmpPercentage = (statistics.snmp_devices / statistics.total_devices) * 100

  return (
    <div className="space-y-6">
      {/* Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Devices</CardTitle>
            <Network className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{statistics.total_devices}</div>
            <p className="text-xs text-muted-foreground">Discovered on {scanInfo.network_range}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Reachable</CardTitle>
            <Wifi className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{statistics.reachable_devices}</div>
            <Progress value={reachabilityPercentage} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-1">{reachabilityPercentage.toFixed(1)}% reachable</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">SNMP Enabled</CardTitle>
            <Database className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">{statistics.snmp_devices}</div>
            <Progress value={snmpPercentage} className="mt-2" />
            <p className="text-xs text-muted-foreground mt-1">{snmpPercentage.toFixed(1)}% with SNMP</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Response</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{statistics.avg_response_time_ms}ms</div>
            <p className="text-xs text-muted-foreground">Average response time</p>
          </CardContent>
        </Card>
      </div>

      {/* Scan Configuration */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Scan Configuration
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <p className="text-sm font-medium">Network Range</p>
              <Badge variant="outline" className="font-mono">
                {scanInfo.network_range}
              </Badge>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium">Scan Type</p>
              <Badge variant="default">{scanInfo.scan_type.toUpperCase()}</Badge>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium">SNMP Communities</p>
              <div className="flex flex-wrap gap-1">
                {scanInfo.snmp_communities.map((community, index) => (
                  <Badge key={index} variant="secondary" className="text-xs">
                    {community}
                  </Badge>
                ))}
              </div>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium">Timeout</p>
              <Badge variant="outline">{scanInfo.timeout}s</Badge>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium">Retries</p>
              <Badge variant="outline">{scanInfo.retries}</Badge>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium">Workers</p>
              <Badge variant="outline">{scanInfo.worker_count}</Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Vendor Distribution */}
        <Card>
          <CardHeader>
            <CardTitle>Vendor Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={vendorData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {vendorData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Scan Method Distribution */}
        <Card>
          <CardHeader>
            <CardTitle>Scan Method Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={methodData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="value" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
