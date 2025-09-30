"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Slider } from "@/components/ui/slider";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Loader2, Play, Settings, Target } from "lucide-react";
import type { ScanResult } from "@/app/page";

interface NetworkScanFormProps {
  onScanStart: () => void;
  onScanComplete: (result: ScanResult) => void;
  onScanError: (error: string) => void;
  isScanning: boolean;
}

export function NetworkScanForm({
  onScanStart,
  onScanComplete,
  onScanError,
  isScanning,
}: NetworkScanFormProps) {
  const [networkRange, setNetworkRange] = useState("192.168.1.0/24");
  const [singleIP, setSingleIP] = useState("192.168.1.1");
  const [communities, setCommunities] = useState("public,private");
  const [timeout, setTimeout] = useState([2]);
  const [retries, setRetries] = useState([1]);
  const [scanType, setScanType] = useState("full");
  const [enablePortScan, setEnablePortScan] = useState(true);
  const [scanMode, setScanMode] = useState<"network" | "single">("network");

  const performNetworkScan = async (scanData: any) => {
    try {
      const response = await fetch(
        "http://localhost:8080/api/v1/network/full-scan",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(scanData),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      return result;
    } catch (error) {
      console.error("Scan failed:", error);
      throw error;
    }
  };

  const performSingleDeviceScan = async (
    ip: string,
    communities: string[],
    enablePortScan: boolean
  ) => {
    try {
      const communityParams = communities
        .map((c) => `community=${c}`)
        .join("&");
      const response = await fetch(
        `http://localhost:8080/api/v1/device/${ip}?${communityParams}&enable_port_scan=${enablePortScan}`
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const deviceResult = await response.json();
      console.log("Device result keys:", Object.keys(deviceResult));
      console.log(
        "Device result structure:",
        JSON.stringify(deviceResult, null, 2)
      );

      // Transform single device result to match ScanResult format
      const result: ScanResult = {
        topology: {
          devices: [deviceResult.device],
          total_count: 1,
          reachable_count: deviceResult.device.is_reachable ? 1 : 0,
          snmp_count: deviceResult.device.scan_method?.includes("SNMP") ? 1 : 0,
          arp_count: deviceResult.device.scan_method?.includes("ARP") ? 1 : 0,
          scan_duration_ms: 1000,
          scan_method: "SINGLE_DEVICE",
        },
        statistics: {
          total_devices: 1,
          reachable_devices: deviceResult.device.is_reachable ? 1 : 0,
          snmp_devices: deviceResult.device.scan_method?.includes("SNMP")
            ? 1
            : 0,
          arp_only_devices: deviceResult.device.scan_method === "ARP" ? 1 : 0,
          devices_with_mac: deviceResult.device.mac_address ? 1 : 0,
          vendor_distribution: { [deviceResult.device.vendor || "Unknown"]: 1 },
          scan_method_distribution: {
            [deviceResult.device.scan_method || "Unknown"]: 1,
          },
          avg_response_time_ms: deviceResult.device.response_time_ms || 0,
        },
        scan_info: {
          scan_type: "single_device",
          network_range: ip,
          snmp_communities: communities,
          timeout: timeout[0],
          retries: retries[0],
          worker_count: 1,
        },
      };

      return result;
    } catch (error) {
      console.error("Single device scan failed:", error);
      throw error;
    }
  };

  const handleScan = async () => {
    if (scanMode === "network" && !networkRange.trim()) {
      onScanError("Network range is required");
      return;
    }

    if (scanMode === "single" && !singleIP.trim()) {
      onScanError("IP address is required");
      return;
    }

    const communityList = communities
      .split(",")
      .map((c) => c.trim())
      .filter((c) => c);

    onScanStart();

    try {
      let result: ScanResult;

      if (scanMode === "single") {
        result = await performSingleDeviceScan(
          singleIP.trim(),
          communityList,
          enablePortScan
        );
      } else {
        const scanData = {
          network_range: networkRange.trim(),
          communities: communityList,
          timeout: timeout[0],
          retries: retries[0],
          scan_type: scanType,
          enable_port_scan: enablePortScan,
        };
        result = await performNetworkScan(scanData);
      }

      console.log("Final single device scan result:", result);
      onScanComplete(result);
    } catch (error) {
      onScanError(error instanceof Error ? error.message : "Scan failed");
    }
  };

  return (
    <Card className="bg-card/80 backdrop-blur-sm border-primary/20">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-foreground">
          <Settings className="h-5 w-5 text-primary" />
          Scan Configuration
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <Tabs
          value={scanMode}
          onValueChange={(value: string) =>
            setScanMode(value as "network" | "single")
          }
          className="w-full"
        >
          <TabsList className="grid w-full grid-cols-2 bg-secondary/50">
            <TabsTrigger
              value="network"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Settings className="h-4 w-4 mr-2" />
              Network Scan
            </TabsTrigger>
            <TabsTrigger
              value="single"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Target className="h-4 w-4 mr-2" />
              Single IP
            </TabsTrigger>
          </TabsList>

          <TabsContent value="network" className="space-y-6 mt-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Network Range */}
              <div className="space-y-2">
                <Label htmlFor="network-range" className="text-foreground">
                  Network Range
                </Label>
                <Input
                  id="network-range"
                  placeholder="192.168.1.0/24"
                  value={networkRange}
                  onChange={(e) => setNetworkRange(e.target.value)}
                  disabled={isScanning}
                  className="bg-input/50 border-border/50"
                />
              </div>

              {/* Scan Type */}
              <div className="space-y-2">
                <Label htmlFor="scan-type" className="text-foreground">
                  Scan Type
                </Label>
                <Select
                  value={scanType}
                  onValueChange={setScanType}
                  disabled={isScanning}
                >
                  <SelectTrigger className="bg-input/50 border-border/50">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-popover border-border/50">
                    <SelectItem value="full">Full Scan</SelectItem>
                    <SelectItem value="snmp">SNMP Only</SelectItem>
                    <SelectItem value="arp">ARP Only</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* SNMP Communities */}
              {scanType !== "arp" && (
                <div className="space-y-2">
                  <Label htmlFor="communities" className="text-foreground">
                    SNMP Communities
                  </Label>
                  <Input
                    id="communities"
                    placeholder="public,private"
                    value={communities}
                    onChange={(e) => setCommunities(e.target.value)}
                    disabled={isScanning}
                    className="bg-input/50 border-border/50"
                  />
                </div>
              )}

              {/* Port Scan Toggle */}
              <div className="space-y-2">
                <Label htmlFor="port-scan" className="text-foreground">
                  Port Scanning
                </Label>
                <div className="flex items-center space-x-2">
                  <Switch
                    id="port-scan"
                    checked={enablePortScan}
                    onCheckedChange={setEnablePortScan}
                    disabled={isScanning}
                  />
                  <Label
                    htmlFor="port-scan"
                    className="text-sm text-muted-foreground"
                  >
                    {enablePortScan ? "Enabled" : "Disabled"}
                  </Label>
                </div>
              </div>
            </div>
          </TabsContent>

          <TabsContent value="single" className="space-y-6 mt-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Single IP Address */}
              <div className="space-y-2">
                <Label htmlFor="single-ip" className="text-foreground">
                  IP Address
                </Label>
                <Input
                  id="single-ip"
                  placeholder="192.168.1.1"
                  value={singleIP}
                  onChange={(e) => setSingleIP(e.target.value)}
                  disabled={isScanning}
                  className="bg-input/50 border-border/50"
                />
              </div>

              {/* SNMP Communities */}
              <div className="space-y-2">
                <Label htmlFor="communities-single" className="text-foreground">
                  SNMP Communities
                </Label>
                <Input
                  id="communities-single"
                  placeholder="public,private"
                  value={communities}
                  onChange={(e) => setCommunities(e.target.value)}
                  disabled={isScanning}
                  className="bg-input/50 border-border/50"
                />
              </div>

              {/* Port Scan Toggle for Single IP */}
              <div className="space-y-2 md:col-span-2">
                <Label htmlFor="port-scan-single" className="text-foreground">
                  Port Scanning
                </Label>
                <div className="flex items-center space-x-2">
                  <Switch
                    id="port-scan-single"
                    checked={enablePortScan}
                    onCheckedChange={setEnablePortScan}
                    disabled={isScanning}
                  />
                  <Label
                    htmlFor="port-scan-single"
                    className="text-sm text-muted-foreground"
                  >
                    {enablePortScan ? "Enabled" : "Disabled"}
                  </Label>
                </div>
              </div>
            </div>
          </TabsContent>
        </Tabs>

        {scanMode === "network" && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Timeout Slider */}
            <div className="space-y-3">
              <Label className="text-foreground">
                Timeout: {timeout[0]} seconds
              </Label>
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
              <Label className="text-foreground">Retries: {retries[0]}</Label>
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
        )}

        {/* Scan Button */}
        <Button
          onClick={handleScan}
          disabled={isScanning}
          className="w-full bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg shadow-primary/20"
          size="lg"
        >
          {isScanning ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              {scanMode === "single"
                ? "Scanning Device..."
                : "Scanning Network..."}
            </>
          ) : (
            <>
              <Play className="mr-2 h-4 w-4" />
              {scanMode === "single"
                ? "Scan Single Device"
                : "Start Network Scan"}
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  );
}
