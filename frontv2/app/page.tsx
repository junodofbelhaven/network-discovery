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
import { Badge } from "@/components/ui/badge"
import { Activity, Shield, Wifi, Zap } from "lucide-react"

export interface Device {
  ip: string
  mac_address?: string
  hostname?: string
  vendor?: string
  description?: string
  uptime?: string
  is_reachable: boolean
  response_time_ms?: number
  scan_method?: string
  open_ports?: Array<{
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
    <div className="dark min-h-screen bg-background">
      <div className="absolute inset-0 bg-gradient-to-br from-background via-background to-primary/5"></div>
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(120,119,198,0.1),transparent_50%)]"></div>

      <div className="relative max-w-7xl mx-auto p-6 space-y-8">
        <div className="text-center space-y-6 py-8">
          <div className="flex justify-center items-center space-x-4 mb-4">
            <div className="relative">
              <Shield className="h-12 w-12 text-primary drop-shadow-[0_0_15px_rgba(120,119,198,0.5)]" />
              <div className="absolute inset-0 pulse-ring border-2 border-primary/60 rounded-full"></div>
            </div>
            <div className="h-8 w-px bg-primary/30"></div>
            <Wifi className="h-10 w-10 text-accent drop-shadow-[0_0_10px_rgba(120,119,198,0.3)]" />
          </div>

          <div className="space-y-3">
            <h1 className="text-5xl font-bold bg-gradient-to-r from-primary via-accent to-chart-1 bg-clip-text text-transparent drop-shadow-sm">
              Network Discovery Tool
            </h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Advanced network reconnaissance and device analysis platform for cybersecurity professionals
            </p>

            <div className="flex justify-center space-x-4 mt-6">
              <Badge variant="secondary" className="px-4 py-2 bg-secondary/80 border border-primary/20">
                <Activity className="h-4 w-4 mr-2" />
                Real-time Scanning
              </Badge>
              <Badge variant="secondary" className="px-4 py-2 bg-secondary/80 border border-primary/20">
                <Zap className="h-4 w-4 mr-2" />
                Port Analysis
              </Badge>
              <Badge variant="secondary" className="px-4 py-2 bg-secondary/80 border border-primary/20">
                <Shield className="h-4 w-4 mr-2" />
                Security Focused
              </Badge>
            </div>
          </div>
        </div>

        <Card className="bg-card/50 backdrop-blur-xl border-2 border-primary/30 shadow-2xl shadow-primary/10">
          <div className="p-8">
            <div className="flex items-center space-x-3 mb-6">
              <div className="h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center border border-primary/30">
                <Wifi className="h-4 w-4 text-primary" />
              </div>
              <h2 className="text-2xl font-semibold text-foreground">Network Scan Configuration</h2>
            </div>
            <NetworkScanForm
              onScanStart={handleScanStart}
              onScanComplete={handleScanComplete}
              onScanError={handleScanError}
              isScanning={isScanning}
            />
          </div>
        </Card>

        {(isScanning || scanResult) && (
          <div className="bg-card/50 backdrop-blur-xl rounded-xl border border-primary/30 shadow-lg shadow-primary/5">
            <ScanStatus isScanning={isScanning} progress={scanProgress} duration={scanDuration} error={error} />
          </div>
        )}

        {scanResult && (
          <div className="space-y-6">
            <Tabs defaultValue="results" className="space-y-6">
              <div className="flex justify-center">
                <TabsList className="bg-card/50 backdrop-blur-xl border border-primary/30 p-1 shadow-lg">
                  <TabsTrigger
                    value="results"
                    className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-lg"
                  >
                    <Activity className="h-4 w-4 mr-2" />
                    Scan Results
                  </TabsTrigger>
                  <TabsTrigger
                    value="statistics"
                    className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground data-[state=active]:shadow-lg"
                  >
                    <Zap className="h-4 w-4 mr-2" />
                    Analytics
                  </TabsTrigger>
                </TabsList>
              </div>

              <TabsContent value="results" className="space-y-6">
                <Card className="bg-card/50 backdrop-blur-xl border-primary/30 shadow-2xl shadow-primary/10">
                  <div className="p-6">
                    <div className="flex justify-between items-center mb-6">
                      <div className="space-y-1">
                        <h2 className="text-2xl font-semibold flex items-center text-foreground">
                          <Shield className="h-6 w-6 mr-3 text-primary drop-shadow-[0_0_10px_rgba(120,119,198,0.3)]" />
                          Discovered Devices
                        </h2>
                        <p className="text-muted-foreground">
                          Found {scanResult.topology.devices?.length || 0} devices on the network
                        </p>
                      </div>
                      <ViewToggle viewMode={viewMode} onViewModeChange={setViewMode} />
                    </div>

                    {viewMode === "table" ? (
                      <DeviceTable devices={scanResult.topology.devices || []} />
                    ) : (
                      <DeviceCards devices={scanResult.topology.devices || []} />
                    )}
                  </div>
                </Card>
              </TabsContent>

              <TabsContent value="statistics">
                <Card className="bg-card/50 backdrop-blur-xl border-primary/30 shadow-2xl shadow-primary/10">
                  <div className="p-6">
                    <div className="flex items-center space-x-3 mb-6">
                      <div className="h-8 w-8 rounded-full bg-accent/20 flex items-center justify-center border border-accent/30">
                        <Zap className="h-4 w-4 text-accent" />
                      </div>
                      <h2 className="text-2xl font-semibold text-foreground">Network Analytics</h2>
                    </div>
                    <StatsDashboard statistics={scanResult.statistics} scanInfo={scanResult.scan_info} />
                  </div>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        )}
      </div>
    </div>
  )
}
