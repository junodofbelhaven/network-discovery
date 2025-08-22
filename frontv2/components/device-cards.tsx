"use client"

import { useState } from "react"
import type { Device } from "@/app/page"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import {
  Search,
  Filter,
  WifiOff,
  Network,
  HardDrive,
  Shield,
  Activity,
  Zap,
  Router,
  Server,
  Smartphone,
  Monitor,
} from "lucide-react"

interface DeviceCardsProps {
  devices: Device[]
}

export function DeviceCards({ devices }: DeviceCardsProps) {
  console.log("First device structure:", devices[0] ? JSON.stringify(devices[0], null, 2) : "No devices")
  const [searchTerm, setSearchTerm] = useState("")
  const [vendorFilter, setVendorFilter] = useState("all")
  const [methodFilter, setMethodFilter] = useState("all")

  const uniqueVendors = Array.from(new Set(devices.map((d) => d.vendor))).filter(Boolean)
  const uniqueMethods = Array.from(new Set(devices.map((d) => d.scan_method))).filter(Boolean)

  const filteredDevices = devices.filter((device) => {
    const matchesSearch = searchTerm === "" || 
      device.ip?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      device.hostname?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      device.vendor?.toLowerCase().includes(searchTerm.toLowerCase())

    const matchesVendor = vendorFilter === "all" || device.vendor === vendorFilter
    const matchesMethod = methodFilter === "all" || device.scan_method === methodFilter

    return matchesSearch && matchesVendor && matchesMethod
  })

  console.log("Filtered devices:", filteredDevices)

  const getStatusColor = (device: Device) => {
    if (!device.is_reachable)
      return "border-destructive/30 bg-gradient-to-br from-destructive/5 via-card to-destructive/10 hover:shadow-destructive/20"

    switch (device.scan_method) {
      case "SNMP":
        return "border-chart-3/30 bg-gradient-to-br from-chart-3/5 via-card to-chart-3/10 hover:shadow-chart-3/20"
      case "ARP":
        return "border-chart-2/30 bg-gradient-to-br from-chart-2/5 via-card to-chart-2/10 hover:shadow-chart-2/20"
      case "COMBINED":
        return "border-chart-1/30 bg-gradient-to-br from-chart-1/5 via-card to-chart-1/10 hover:shadow-chart-1/20"
      default:
        return "border-primary/30 bg-gradient-to-br from-primary/5 via-card to-primary/10 hover:shadow-primary/20"
    }
  }

  const getStatusBadge = (device: Device) => {
    if (!device.is_reachable) {
      return (
        <Badge className="bg-destructive hover:bg-destructive/90 text-destructive-foreground shadow-lg border-0">
          Offline
        </Badge>
      )
    }

    switch (device.scan_method) {
      case "SNMP":
        return <Badge className="bg-chart-3 hover:bg-chart-3/90 text-white shadow-lg border-0">SNMP</Badge>
      case "ARP":
        return <Badge className="bg-chart-2 hover:bg-chart-2/90 text-white shadow-lg border-0">ARP</Badge>
      case "COMBINED":
        return <Badge className="bg-chart-1 hover:bg-chart-1/90 text-white shadow-lg border-0">Full Scan</Badge>
      default:
        return (
          <Badge className="bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg border-0">
            {device.scan_method}
          </Badge>
        )
    }
  }

  const getDeviceIcon = (device: Device) => {
    const vendor = device.vendor?.toLowerCase() || ""
    const hostname = device.hostname?.toLowerCase() || ""

    if (vendor.includes("cisco") || vendor.includes("router") || hostname.includes("router")) {
      return <Router className="h-6 w-6 text-primary" />
    }
    if (vendor.includes("server") || hostname.includes("server")) {
      return <Server className="h-6 w-6 text-chart-3" />
    }
    if (vendor.includes("apple") || vendor.includes("samsung") || hostname.includes("phone")) {
      return <Smartphone className="h-6 w-6 text-chart-2" />
    }
    if (vendor.includes("dell") || vendor.includes("hp") || hostname.includes("pc")) {
      return <Monitor className="h-6 w-6 text-chart-1" />
    }
    return <Network className="h-6 w-6 text-muted-foreground" />
  }

  return (
    <div className="space-y-8">
      <Card className="glass-effect border-primary/20 shadow-xl">
        <CardHeader className="pb-6">
          <CardTitle className="flex items-center gap-3 text-card-foreground text-xl">
            <div className="p-2 rounded-lg bg-primary/10">
              <Filter className="h-5 w-5 text-primary" />
            </div>
            Device Discovery Results
          </CardTitle>
          <div className="flex flex-col sm:flex-row gap-4 mt-4">
            <div className="relative flex-1">
              <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-muted-foreground" />
              <Input
                placeholder="Search by IP, hostname, or vendor..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-12 h-12 bg-input/50 border-border/50 focus:border-primary/50 text-lg"
              />
            </div>
            <Select value={vendorFilter} onValueChange={setVendorFilter}>
              <SelectTrigger className="w-full sm:w-56 h-12 bg-input/50 border-border/50">
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
              <SelectTrigger className="w-full sm:w-56 h-12 bg-input/50 border-border/50">
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

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
        {filteredDevices.map((device, index) => (
          <Card
            key={device.ip || `device-${index}`}
            className={`
              group relative overflow-hidden transition-all duration-500 
              hover:shadow-2xl hover:-translate-y-3 hover:scale-[1.03]
              border-2 shadow-xl backdrop-blur-md
              ${getStatusColor(device)}
            `}
          >
            <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-accent/5 opacity-0 group-hover:opacity-100 transition-opacity duration-700" />
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-primary via-accent to-primary opacity-60" />

            <CardHeader className="relative pb-4 space-y-4">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                  <div className="p-3 rounded-xl bg-card/80 shadow-lg border border-border/30 group-hover:shadow-xl transition-all duration-300">
                    {getDeviceIcon(device)}
                  </div>
                  <div className="space-y-2">
                    <CardTitle className="text-xl font-mono font-bold text-card-foreground tracking-tight">
                      {device.ip}
                    </CardTitle>
                    {device.hostname && (
                      <p className="text-sm font-semibold text-muted-foreground truncate max-w-[200px] bg-muted/20 px-2 py-1 rounded-md">
                        {device.hostname}
                      </p>
                    )}
                  </div>
                </div>
                <div className="flex flex-col items-end gap-3">
                  {device.is_reachable ? (
                    <div className="flex items-center gap-2 px-3 py-2 rounded-full bg-chart-1/20 border border-chart-1/30">
                      <div className="w-3 h-3 bg-chart-1 rounded-full animate-pulse shadow-lg" />
                      <span className="text-sm font-bold text-chart-1">Online</span>
                    </div>
                  ) : (
                    <div className="flex items-center gap-2 px-3 py-2 rounded-full bg-destructive/20 border border-destructive/30">
                      <WifiOff className="h-4 w-4 text-destructive" />
                      <span className="text-sm font-bold text-destructive">Offline</span>
                    </div>
                  )}
                  {getStatusBadge(device)}
                </div>
              </div>
            </CardHeader>

            <CardContent className="relative space-y-5">
              <div className="grid grid-cols-1 gap-4">
                <div className="group/item flex items-center gap-4 p-4 glass-effect rounded-xl border border-border/30 hover:border-primary/30 transition-all duration-300 hover:shadow-lg">
                  <div className="flex-shrink-0 p-3 rounded-xl bg-primary/10 group-hover/item:bg-primary/20 transition-all duration-300">
                    <HardDrive className="h-5 w-5 text-primary" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-1">Vendor</p>
                    <p className="text-base font-bold text-card-foreground truncate">{device.vendor || "Unknown"}</p>
                  </div>
                </div>

                {device.mac_address && (
                  <div className="group/item flex items-center gap-4 p-4 glass-effect rounded-xl border border-border/30 hover:border-accent/30 transition-all duration-300 hover:shadow-lg">
                    <div className="flex-shrink-0 p-3 rounded-xl bg-accent/10 group-hover/item:bg-accent/20 transition-all duration-300">
                      <Network className="h-5 w-5 text-accent" />
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-1">
                        MAC Address
                      </p>
                      <p className="text-sm font-mono font-bold text-card-foreground tracking-wider bg-muted/20 px-2 py-1 rounded">
                        {device.mac_address}
                      </p>
                    </div>
                  </div>
                )}

                <div className="group/item flex items-center gap-4 p-4 glass-effect rounded-xl border border-border/30 hover:border-chart-2/30 transition-all duration-300 hover:shadow-lg">
                  <div className="flex-shrink-0 p-3 rounded-xl bg-chart-2/10 group-hover/item:bg-chart-2/20 transition-all duration-300">
                    <Zap className="h-5 w-5 text-chart-2" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-bold text-muted-foreground uppercase tracking-widest mb-1">
                      Response Time
                    </p>
                    <p className="text-base font-bold text-card-foreground">
                      {device.response_time_ms}
                      <span className="text-sm text-muted-foreground ml-1">ms</span>
                    </p>
                  </div>
                </div>
              </div>

              {device.open_ports && device.open_ports.length > 0 && (
                <div className="space-y-4 p-5 glass-effect rounded-xl border border-border/30">
                  <div className="flex items-center gap-3">
                    <div className="p-2 rounded-lg bg-chart-4/10">
                      <Shield className="h-5 w-5 text-chart-4" />
                    </div>
                    <p className="text-base font-bold text-card-foreground">
                      Open Ports
                      <span className="ml-2 px-3 py-1 text-sm bg-chart-4/20 text-chart-4 rounded-full font-mono">
                        {device.open_ports.length}
                      </span>
                    </p>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {device.open_ports.slice(0, 6).map((port) => (
                      <Badge
                        key={port.port}
                        className="text-sm font-mono bg-muted/30 hover:bg-muted/50 text-card-foreground border border-border/50 shadow-md px-3 py-1"
                      >
                        {port.port}/{port.protocol}
                      </Badge>
                    ))}
                    {device.open_ports.length > 6 && (
                      <Badge className="text-sm bg-primary/20 hover:bg-primary/30 text-primary border border-primary/30 px-3 py-1">
                        +{device.open_ports.length - 6} more
                      </Badge>
                    )}
                  </div>
                </div>
              )}

              {device.uptime && (
                <div className="pt-4 border-t border-border/30">
                  <div className="flex items-center gap-3 text-sm text-muted-foreground">
                    <div className="p-1.5 rounded-lg bg-muted/20">
                      <Activity className="h-4 w-4" />
                    </div>
                    <span className="font-semibold">Uptime:</span>
                    <span className="font-mono font-bold text-card-foreground">{device.uptime}</span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {filteredDevices.length === 0 && (
        <Card className="glass-effect border-primary/20 shadow-xl">
          <CardContent className="text-center py-16">
            <div className="space-y-4">
              <div className="mx-auto w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                <Search className="h-8 w-8 text-primary" />
              </div>
              <p className="text-xl text-card-foreground font-semibold">No devices found</p>
              <p className="text-muted-foreground">Try adjusting your search criteria or scan parameters</p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
