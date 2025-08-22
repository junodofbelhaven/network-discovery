"use client"

import { useState } from "react"
import type { Device } from "@/app/page"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Search, Filter, Wifi, WifiOff } from "lucide-react"

interface DeviceTableProps {
  devices: Device[]
}

export function DeviceTable({ devices }: DeviceTableProps) {
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
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Filter className="h-5 w-5" />
          Device List
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
      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Status</TableHead>
                <TableHead>IP Address</TableHead>
                <TableHead>Hostname</TableHead>
                <TableHead>MAC Address</TableHead>
                <TableHead>Vendor</TableHead>
                <TableHead>Response Time</TableHead>
                <TableHead>Open Ports</TableHead>
                <TableHead>Method</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredDevices.map((device) => (
                <TableRow key={device.ip}>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      {device.is_reachable ? (
                        <Wifi className="h-4 w-4 text-green-500" />
                      ) : (
                        <WifiOff className="h-4 w-4 text-red-500" />
                      )}
                      {getStatusBadge(device)}
                    </div>
                  </TableCell>
                  <TableCell className="font-mono">{device.ip}</TableCell>
                  <TableCell>{device.hostname || "-"}</TableCell>
                  <TableCell className="font-mono text-sm">{device.mac_address || "-"}</TableCell>
                  <TableCell>{device.vendor || "Unknown"}</TableCell>
                  <TableCell>{device.response_time_ms}ms</TableCell>
                  <TableCell>
                    {device.open_ports && device.open_ports.length > 0 ? (
                      <div className="flex flex-wrap gap-1">
                        {device.open_ports.slice(0, 3).map((port) => (
                          <Badge key={port.port} variant="outline" className="text-xs">
                            {port.port}/{port.protocol}
                          </Badge>
                        ))}
                        {device.open_ports.length > 3 && (
                          <Badge variant="outline" className="text-xs">
                            +{device.open_ports.length - 3}
                          </Badge>
                        )}
                      </div>
                    ) : (
                      "-"
                    )}
                  </TableCell>
                  <TableCell>{device.scan_method}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
        {filteredDevices.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">No devices found matching your criteria</div>
        )}
      </CardContent>
    </Card>
  )
}
