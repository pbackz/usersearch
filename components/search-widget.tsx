"use client"

import { useState } from "react"
import { Search, Loader2, ExternalLink } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card } from "@/components/ui/card"

interface SearchResult {
  platform: string;
  url: string;
}

export default function SearchWidget() {
  const [query, setQuery] = useState("")
  const [results, setResults] = useState<SearchResult[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")

  const handleSearch = async () => {
    if (!query.trim()) return
    
    setLoading(true)
    setError("")
    try {
      const response = await fetch("/api/search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ query }),
      })
      
      if (!response.ok) {
        throw new Error("Search failed")
      }
      
      const data = await response.json()
      setResults(data)
    } catch (error) {
      console.error("Search failed:", error)
      setError("Search failed. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <Input
          placeholder="Search GitHub repositories..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="flex-1"
          onKeyDown={(e) => e.key === "Enter" && handleSearch()}
        />
        <Button onClick={handleSearch} disabled={loading}>
          {loading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Search className="h-4 w-4" />
          )}
          <span className="ml-2">Search</span>
        </Button>
      </div>

      {error && (
        <div className="text-destructive text-sm">{error}</div>
      )}

      {results.length > 0 && (
        <Card className="divide-y divide-border">
          {results.map((result, index) => (
            <div key={index} className="p-4 grid grid-cols-[200px,1fr,auto] gap-4 items-center">
              <div className="font-medium">{result.platform}</div>
              <a 
                href={result.url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-muted-foreground hover:text-primary truncate"
              >
                {result.url}
              </a>
              <Button
                variant="ghost"
                size="icon"
                asChild
              >
                <a
                  href={result.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-muted-foreground hover:text-primary"
                >
                  <ExternalLink className="h-4 w-4" />
                </a>
              </Button>
            </div>
          ))}
        </Card>
      )}
    </div>
  )
}