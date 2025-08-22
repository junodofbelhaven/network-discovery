"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Slider } from "@/components/ui/slider"
import { Switch } from "@/components/ui/switch"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2, Play, Settings } from "lucide-react"
import type { ScanResult } from "@/app/page"

interface NetworkScanFormProps {
  onScanStart: () => void
  onScanComplete: (result: ScanResult) => void
  onScanError: (error: string) => void
  isScanning: boolean
}

export function NetworkScanForm({ onScanStart, onScanComplete, onScanError, isScanning }: NetworkScanFormProps) {
  const [networkRange, setNetworkRange] = useState("192.168.1.0/24")
  const [communities, setCommunities] = useState("public,private")
  const [timeout, setTimeout] = useState([2])
  const [retries, setRetries] = useState([1])
  const [scanType, setScanType] = useState("full")
  const [enablePortScan, setEnablePortScan] = useState(true)

  const performNetworkScan = async (scanData: any) => {
    try {
      const response = await fetch("http://localhost:8080/api/v1/network/full-scan", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(scanData),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      return result
    } catch (error) {
      console.error("Scan failed:", error)
      throw error
    }
  }

  const handleScan = async () => {
    if (!networkRange.trim()) {
      onScanError("Network range is required")
      return
    }

    const scanData = {
      network_range: networkRange.trim(),
      communities: communities
        .split(",")
        .map((c) => c.trim())
        .filter((c) => c),
      timeout: timeout[0],
      retries: retries[0],
      scan_type: scanType,
      enable_port_scan: enablePortScan,
    }

    onScanStart()

    try {
      const result = await performNetworkScan(scanData)
      onScanComplete(result)
    } catch (error) {
      onScanError(error instanceof Error ? error.message : "Scan failed")
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings className="h-5 w-5" />
          Scan Configuration
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Network Range */}
          <div className="space-y-2">
            <Label htmlFor="network-range">Network Range</Label>
            <Input
              id="network-range"
              placeholder="192.168.1.0/24"
              value={networkRange}
              onChange={(e) => setNetworkRange(e.target.value)}
              disabled={isScanning}
            />
          </div>

          {/* SNMP Communities */}
          <div className="space-y-2">
            <Label htmlFor="communities">SNMP Communities</Label>
            <Input
              id="communities"
              placeholder="public,private"
              value={communities}
              onChange={(e) => setCommunities(e.target.value)}
              disabled={isScanning}
            />
          </div>

          {/* Scan Type */}
          <div className="space-y-2">
            <Label htmlFor="scan-type">Scan Type</Label>
            <Select value={scanType} onValueChange={setScanType} disabled={isScanning}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="full">Full Scan</SelectItem>
                <SelectItem value="snmp">SNMP Only</SelectItem>
                <SelectItem value="arp">ARP Only</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Port Scan Toggle */}
          <div className="space-y-2">
            <Label htmlFor="port-scan">Port Scanning</Label>
            <div className="flex items-center space-x-2">
              <Switch
                id="port-scan"
                checked={enablePortScan}
                onCheckedChange={setEnablePortScan}
                disabled={isScanning}
              />
              <Label htmlFor="port-scan" className="text-sm text-muted-foreground">
                {enablePortScan ? "Enabled" : "Disabled"}
              </Label>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Timeout Slider */}
          <div className="space-y-3">
            <Label>Timeout: {timeout[0]} seconds</Label>
            <Slider
              value={timeout}
              onValueChange={setTimeout}
              max={10}
              min={1}
              step={1}
              disabled={isScanning}
              className="w-full"
            />
          </div>

          {/* Retries Slider */}
          <div className="space-y-3">
            <Label>Retries: {retries[0]}</Label>
            <Slider
              value={retries}
              onValueChange={setRetries}
              max={3}
              min={0}
              step={1}
              disabled={isScanning}
              className="w-full"
            />
          </div>
        </div>

        {/* Scan Button */}
        <Button onClick={handleScan} disabled={isScanning} className="w-full" size="lg">
          {isScanning ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Scanning Network...
            </>
          ) : (
            <>
              <Play className="mr-2 h-4 w-4" />
              Start Network Scan
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  )
}
