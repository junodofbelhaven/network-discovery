"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Loader2, CheckCircle, XCircle, Clock } from "lucide-react"

interface ScanStatusProps {
  isScanning: boolean
  progress: number
  duration: number
  error: string | null
}

export function ScanStatus({ isScanning, progress, duration, error }: ScanStatusProps) {
  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, "0")}`
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <XCircle className="h-4 w-4" />
        <AlertDescription>Scan failed: {error}</AlertDescription>
      </Alert>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          {isScanning ? (
            <>
              <Loader2 className="h-5 w-5 animate-spin" />
              Scanning Network...
            </>
          ) : (
            <>
              <CheckCircle className="h-5 w-5 text-green-500" />
              Scan Complete
            </>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-sm font-medium">Progress</span>
            <Badge variant={isScanning ? "default" : "secondary"}>{progress}%</Badge>
          </div>
          <Progress value={progress} className="w-full" />
        </div>

        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Clock className="h-4 w-4" />
          <span>Duration: {formatDuration(duration)}</span>
        </div>
      </CardContent>
    </Card>
  )
}
