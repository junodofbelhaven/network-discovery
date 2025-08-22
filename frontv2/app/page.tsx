"use client"

import { useState, useEffect } from "react"
import { NetworkScanForm } from "@/components/network-scan-form"
import { DeviceTable } from "@/components/device-table"
import { DeviceCards } from "@/components/device-cards"
import { StatsDashboard } from "@/components/stats-dashboard"
import { ScanStatus } from "@/components/scan-status"
import { ViewToggle } from "@/components/view-toggle"
import { Card } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

export interface Device {
  ip: string
  mac_address: string
  hostname: string
  vendor: string
  description: string
  uptime: string
  is_reachable: boolean
  response_time_ms: number
  scan_method: string
  open_ports: Array<{
    port: number
    protocol: string
    service: string
    state: string
  }>
}

export interface ScanResult {
  topology: {
    devices: Device[]
    total_count: number
    reachable_count: number
    snmp_count: number
    arp_count: number
    scan_duration_ms: number
    scan_method: string
  }
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
  scan_info: {
    scan_type: string
    network_range: string
    snmp_communities: string[]
    timeout: number
    retries: number
    worker_count: number
  }
}

export default function NetworkDiscoveryTool() {
  const [scanResult, setScanResult] = useState<ScanResult | null>(null)
  const [isScanning, setIsScanning] = useState(false)
  const [scanProgress, setScanProgress] = useState(0)
  const [scanDuration, setScanDuration] = useState(0)
  const [viewMode, setViewMode] = useState<"table" | "cards">("table")
  const [error, setError] = useState<string | null>(null)

  const handleScanComplete = (result: ScanResult) => {
    setScanResult(result)
    setIsScanning(false)
    setScanProgress(100)
    setError(null)
  }

  const handleScanStart = () => {
    setIsScanning(true)
    setScanProgress(0)
    setScanDuration(0)
    setError(null)
  }

  const handleScanError = (errorMessage: string) => {
    setError(errorMessage)
    setIsScanning(false)
    setScanProgress(0)
  }

  useEffect(() => {
    let interval: NodeJS.Timeout
    if (isScanning) {
      interval = setInterval(() => {
        setScanDuration((prev) => prev + 1)
        setScanProgress((prev) => Math.min(prev + 2, 95))
      }, 1000)
    }
    return () => clearInterval(interval)
  }, [isScanning])

  return (
    <div className="min-h-screen bg-background p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold text-foreground">Network Discovery Tool</h1>
          <p className="text-muted-foreground">Discover and analyze devices on your network</p>
        </div>

        {/* Scan Form */}
        <Card className="p-6">
          <NetworkScanForm
            onScanStart={handleScanStart}
            onScanComplete={handleScanComplete}
            onScanError={handleScanError}
            isScanning={isScanning}
          />
        </Card>

        {/* Scan Status */}
        {(isScanning || scanResult) && (
          <ScanStatus isScanning={isScanning} progress={scanProgress} duration={scanDuration} error={error} />
        )}

        {/* Results */}
        {scanResult && (
          <Tabs defaultValue="results" className="space-y-6">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="results">Scan Results</TabsTrigger>
              <TabsTrigger value="statistics">Statistics</TabsTrigger>
            </TabsList>

            <TabsContent value="results" className="space-y-4">
              <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Discovered Devices ({scanResult.topology.devices.length})</h2>
                <ViewToggle viewMode={viewMode} onViewModeChange={setViewMode} />
              </div>

              {viewMode === "table" ? (
                <DeviceTable devices={scanResult.topology.devices} />
              ) : (
                <DeviceCards devices={scanResult.topology.devices} />
              )}
            </TabsContent>

            <TabsContent value="statistics">
              <StatsDashboard statistics={scanResult.statistics} scanInfo={scanResult.scan_info} />
            </TabsContent>
          </Tabs>
        )}
      </div>
    </div>
  )
}
