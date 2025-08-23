"use client";

import { useState } from "react";
import type { Device } from "@/app/page";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Search,
  Filter,
  Wifi,
  WifiOff,
  ChevronUp,
  ChevronDown,
  Network,
  ChevronRight,
  ChevronLeft,
} from "lucide-react";
import { match } from "assert";

interface DeviceTableProps {
  devices: Device[];
}

type SortField = "vendor" | "response_time" | "open_ports" | "hostname";
type SortDirection = "asc" | "desc";

export function DeviceTable({ devices }: DeviceTableProps) {
  const [searchTerm, setSearchTerm] = useState("");
  const [vendorFilter, setVendorFilter] = useState("all");
  const [methodFilter, setMethodFilter] = useState("all");
  const [portFilter, setPortFilter] = useState("all");
  const [sortField, setSortField] = useState<SortField | null>(null);
  const [sortDirection, setSortDirection] = useState<SortDirection>("asc");
  const [showOpenPortsOnly, setShowOpenPortsOnly] = useState(false);
  const [expandedPorts, setExpandedPorts] = useState<Set<string>>(new Set());

  const uniqueVendors = Array.from(
    new Set(devices.map((d) => d.vendor))
  ).filter(Boolean) as string[];
  const uniqueMethods = Array.from(
    new Set(devices.map((d) => d.scan_method))
  ).filter(Boolean) as string[];
  const uniquePorts = Array.from(
    new Set(
      devices
        .flatMap((d) => d.open_ports || [])
        .map((p) => `${p.port}/${p.protocol}`)
    )
  ).sort((a, b) => {
    const aPort = Number.parseInt(a.split("/")[0]);
    const bPort = Number.parseInt(b.split("/")[0]);
    return aPort - bPort;
  });

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      // Set default sort direction based on field
      if (field === "open_ports") {
        setSortDirection("desc"); // Open ports: large to small by default
      } else {
        setSortDirection("asc"); // Others: alphabetical/small to large by default
      }
    }
  };

  const togglePortExpansion = (deviceIp: string) => {
    const newExpanded = new Set(expandedPorts);
    if (newExpanded.has(deviceIp)) {
      newExpanded.delete(deviceIp);
    } else {
      newExpanded.add(deviceIp);
    }
    setExpandedPorts(newExpanded);
  };

  const filteredAndSortedDevices = devices
    .filter((device) => {
      const matchesSearch =
        device.ip.toLowerCase().includes(searchTerm.toLowerCase()) ||
        device.hostname?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        device.vendor?.toLowerCase().includes(searchTerm.toLowerCase());

      const matchesVendor =
        vendorFilter === "all" || device.vendor === vendorFilter;
      const matchesMethod =
        methodFilter === "all" || device.scan_method === methodFilter;
      const matchesPort =
        portFilter === "all" ||
        (device.open_ports &&
          device.open_ports.some(
            (p) => `${p.port}/${p.protocol}` === portFilter
          ));
      const matchesOpenPorts =
        !showOpenPortsOnly ||
        (device.open_ports && device.open_ports.length > 0);

      return (
        matchesSearch &&
        matchesVendor &&
        matchesMethod &&
        matchesOpenPorts &&
        matchesPort
      );
    })
    .sort((a, b) => {
      if (!sortField) return 0;

      let comparison = 0;

      switch (sortField) {
        case "vendor":
          comparison = (a.vendor || "Unknown").localeCompare(
            b.vendor || "Unknown"
          );
          break;
        case "hostname":
          comparison = (a.hostname || "").localeCompare(b.hostname || "");
          break;
        case "response_time":
          comparison = (a.response_time_ms || 0) - (b.response_time_ms || 0);
          break;
        case "open_ports":
          const aPortCount = a.open_ports ? a.open_ports.length : 0;
          const bPortCount = b.open_ports ? b.open_ports.length : 0;
          comparison = aPortCount - bPortCount;
          break;
      }

      return sortDirection === "asc" ? comparison : -comparison;
    });

  const getSortIcon = (field: SortField) => {
    if (sortField !== field) return null;
    return sortDirection === "asc" ? (
      <ChevronUp className="h-4 w-4 ml-1" />
    ) : (
      <ChevronDown className="h-4 w-4 ml-1" />
    );
  };

  const getStatusBadge = (device: Device) => {
    if (!device.is_reachable) {
      return <Badge variant="destructive">Unreachable</Badge>;
    }

    switch (device.scan_method) {
      case "SNMP":
        return <Badge variant="default">SNMP</Badge>;
      case "ARP":
        return <Badge variant="secondary">ARP Only</Badge>;
      case "COMBINED":
        return <Badge variant="outline">Combined</Badge>;
      default:
        return <Badge variant="outline">{device.scan_method}</Badge>;
    }
  };

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
          <Select value={portFilter} onValueChange={setPortFilter}>
            <SelectTrigger className="w-full sm:w-48">
              <SelectValue placeholder="Filter by port" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Ports</SelectItem>
              {uniquePorts.map((port) => (
                <SelectItem key={port} value={port}>
                  Port {port}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button
            variant={showOpenPortsOnly ? "default" : "outline"}
            size="sm"
            onClick={() => setShowOpenPortsOnly(!showOpenPortsOnly)}
            className="flex items-center gap-2 whitespace-nowrap"
          >
            <Network className="h-4 w-4" />
            Open Ports Only
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Status</TableHead>
                <TableHead>IP Address</TableHead>
                <TableHead
                  className="cursor-pointer hover:bg-muted/50 select-none"
                  onClick={() => handleSort("hostname")}
                >
                  <div className="flex items-center">
                    Hostname
                    {getSortIcon("hostname")}
                  </div>
                </TableHead>
                <TableHead>MAC Address</TableHead>
                <TableHead
                  className="cursor-pointer hover:bg-muted/50 select-none"
                  onClick={() => handleSort("vendor")}
                >
                  <div className="flex items-center">
                    Vendor
                    {getSortIcon("vendor")}
                  </div>
                </TableHead>
                <TableHead
                  className="cursor-pointer hover:bg-muted/50 select-none"
                  onClick={() => handleSort("response_time")}
                >
                  <div className="flex items-center">
                    Response Time
                    {getSortIcon("response_time")}
                  </div>
                </TableHead>
                <TableHead
                  className="cursor-pointer hover:bg-muted/50 select-none"
                  onClick={() => handleSort("open_ports")}
                >
                  <div className="flex items-center">
                    Open Ports
                    {getSortIcon("open_ports")}
                  </div>
                </TableHead>
                <TableHead>Method</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredAndSortedDevices.map((device) => (
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
                  <TableCell className="font-mono text-sm">
                    {device.mac_address || "-"}
                  </TableCell>
                  <TableCell>{device.vendor || "Unknown"}</TableCell>
                  <TableCell>{device.response_time_ms}ms</TableCell>
                  <TableCell>
                    {device.open_ports && device.open_ports.length > 0 ? (
                      <div className="flex flex-wrap gap-1">
                        {(expandedPorts.has(device.ip)
                          ? device.open_ports
                          : device.open_ports.slice(0, 3)
                        ).map((port) => (
                          <Badge
                            key={port.port}
                            variant="outline"
                            className="text-xs"
                          >
                            {port.port}/{port.protocol}
                          </Badge>
                        ))}
                        {device.open_ports.length > 3 && (
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-6 px-2 text-xs"
                            onClick={() => togglePortExpansion(device.ip)}
                          >
                            {expandedPorts.has(device.ip) ? (
                              <>
                                <ChevronLeft className="h-3 w-3 mr-1" />
                                Show Less
                              </>
                            ) : (
                              <>
                                <ChevronRight className="h-3 w-3 mr-1" />+
                                {device.open_ports.length - 3} more
                              </>
                            )}
                          </Button>
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
        {filteredAndSortedDevices.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">
            No devices found matching your criteria
          </div>
        )}
      </CardContent>
    </Card>
  );
}
