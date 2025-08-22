"use client"

import { Button } from "@/components/ui/button"
import { Table, Grid3X3 } from "lucide-react"

interface ViewToggleProps {
  viewMode: "table" | "cards"
  onViewModeChange: (mode: "table" | "cards") => void
}

export function ViewToggle({ viewMode, onViewModeChange }: ViewToggleProps) {
  return (
    <div className="flex items-center gap-1 border rounded-lg p-1">
      <Button
        variant={viewMode === "table" ? "default" : "ghost"}
        size="sm"
        onClick={() => onViewModeChange("table")}
        className="h-8 px-3"
      >
        <Table className="h-4 w-4 mr-1" />
        Table
      </Button>
      <Button
        variant={viewMode === "cards" ? "default" : "ghost"}
        size="sm"
        onClick={() => onViewModeChange("cards")}
        className="h-8 px-3"
      >
        <Grid3X3 className="h-4 w-4 mr-1" />
        Cards
      </Button>
    </div>
  )
}
