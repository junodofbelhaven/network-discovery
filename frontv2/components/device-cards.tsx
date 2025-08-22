"use client"

import { useState } from "react"
import type { Device } from "@/app/page"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Search, Filter, Wifi, WifiOff, Clock, Network, HardDrive } from "lucide-react"

interface DeviceCardsProps {
  devices: Device[]
}

export function DeviceCards({ devices }: DeviceCardsProps) {
  const [searchTerm, setSearchTerm] = useState("")
  const [vendorFilter, setVendorFilter] = useState("all")
  const [methodFilter, setMethodFilter] = useState("all")

  const uniqueVendors = Array.from(new Set(devices.map((d) => d.vendor))).filter(Boolean)
  const uniqueMethods = Array.from(new Set(devices.map((d) => d.scan_method))).filter(Boolean)

  const filteredDevices = devices.filter((device) => {
    const matchesSearch =
      device.ip.toLowerCase().includes(searchTerm.toLowerCase()) ||
      device.hostname.toLowerCase().includes(searchTerm.toLowerCase()) ||
      device.vendor.toLowerCase().includes(searchTerm.toLowerCase())

    const matchesVendor = vendorFilter === "all" || device.vendor === vendorFilter
    const matchesMethod = methodFilter === "all" || device.scan_method === methodFilter

    return matchesSearch && matchesVendor && matchesMethod
  })

  const getStatusColor = (device: Device) => {
    if (!device.is_reachable) return "border-red-200 bg-red-50"

    switch (device.scan_method) {
      case "SNMP":
        return "border-blue-200 bg-blue-50"
      case "ARP":
        return "border-orange-200 bg-orange-50"
      case "COMBINED":
        return "border-green-200 bg-green-50"
      default:
        return "border-gray-200 bg-gray-50"
    }
  }

  const getStatusBadge = (device: Device) => {
    if (!device.is_reachable) {
      return <Badge variant="destructive">Unreachable</Badge>
    }

    switch (device.scan_method) {
      case "SNMP":
        return <Badge variant="default">SNMP</Badge>
      case "ARP":
        return <Badge variant="secondary">ARP Only</Badge>
      case "COMBINED":
        return <Badge variant="outline">Combined</Badge>
      default:
        return <Badge variant="outline">{device.scan_method}</Badge>
    }
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Device Cards
          </CardTitle>
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by IP, hostname, or vendor..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10"
              />
            </div>
            <Select value={vendorFilter} onValueChange={setVendorFilter}>
              <SelectTrigger className="w-full sm:w-48">
                <SelectValue placeholder="Filter by vendor" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Vendors</SelectItem>
                {uniqueVendors.map((vendor) => (
                  <SelectItem key={vendor} value={vendor}>
                    {vendor}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={methodFilter} onValueChange={setMethodFilter}>
              <SelectTrigger className="w-full sm:w-48">
                <SelectValue placeholder="Filter by method" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Methods</SelectItem>
                {uniqueMethods.map((method) => (
                  <SelectItem key={method} value={method}>
                    {method}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
      </Card>

      {/* Device Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {filteredDevices.map((device) => (
          <Card key={device.ip} className={`transition-all hover:shadow-md ${getStatusColor(device)}`}>
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg font-mono">{device.ip}</CardTitle>
                <div className="flex items-center gap-2">
                  {device.is_reachable ? (
                    <Wifi className="h-4 w-4 text-green-500" />
                  ) : (
                    <WifiOff className="h-4 w-4 text-red-500" />
                  )}
                  {getStatusBadge(device)}
                </div>
              </div>
              {device.hostname && <p className="text-sm text-muted-foreground">{device.hostname}</p>}
            </CardHeader>
            <CardContent className="space-y-3">
              {/* Device Info */}
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm">
                  <HardDrive className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Vendor:</span>
                  <span>{device.vendor || "Unknown"}</span>
                </div>

                {device.mac_address && (
                  <div className="flex items-center gap-2 text-sm">
                    <Network className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">MAC:</span>
                    <span className="font-mono text-xs">{device.mac_address}</span>
                  </div>
                )}

                <div className="flex items-center gap-2 text-sm">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Response:</span>
                  <span>{device.response_time_ms}ms</span>
                </div>
              </div>

              {/* Open Ports */}
              {device.open_ports && device.open_ports.length > 0 && (
                <div className="space-y-2">
                  <p className="text-sm font-medium">Open Ports:</p>
                  <div className="flex flex-wrap gap-1">
                    {device.open_ports.slice(0, 6).map((port) => (
                      <Badge key={port.port} variant="outline" className="text-xs">
                        {port.port}/{port.protocol}
                      </Badge>
                    ))}
                    {device.open_ports.length > 6 && (
                      <Badge variant="outline" className="text-xs">
                        +{device.open_ports.length - 6} more
                      </Badge>
                    )}
                  </div>
                </div>
              )}

              {/* Uptime */}
              {device.uptime && <div className="text-xs text-muted-foreground">Uptime: {device.uptime}</div>}
            </CardContent>
          </Card>
        ))}
      </div>

      {filteredDevices.length === 0 && (
        <Card>
          <CardContent className="text-center py-8 text-muted-foreground">
            No devices found matching your criteria
          </CardContent>
        </Card>
      )}
    </div>
  )
}
