import { Terminal } from "lucide-react"
import SearchWidget from "@/components/search-widget"

export default function Home() {
  return (
    <main className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="container mx-auto px-4 pt-20 pb-16 text-center">
        <div className="flex items-center justify-center mb-6">
          <Terminal className="h-12 w-12 text-primary" />
        </div>
        <h1 className="text-4xl font-bold tracking-tight mb-4">
          Search username
        </h1>
        <p className="text-muted-foreground text-lg max-w-2xl mx-auto mb-12">
          OSINT tool for searching people's digital footprint. Fast, efficient, and completely free.
        </p>
        
        {/* Search Widget */}
        <div className="max-w-3xl mx-auto">
          <SearchWidget />
        </div>
      </section>

      {/* Features */}
      <section className="container mx-auto px-4 py-16">
        <div className="grid md:grid-cols-3 gap-8">
          <div className="p-6 bg-card rounded-lg">
            <h3 className="text-xl font-semibold mb-3">Lightning Fast</h3>
            <p className="text-muted-foreground">
              Powered by NextJs and Go, search results appear quickly
            </p>
          </div>
          <div className="p-6 bg-card rounded-lg">
            <h3 className="text-xl font-semibold mb-3">Beautiful UI</h3>
            <p className="text-muted-foreground">
              Clean and modern interface for a great user experience
            </p>
          </div>
          <div className="p-6 bg-card rounded-lg">
            <h3 className="text-xl font-semibold mb-3">Open Source</h3>
            <p className="text-muted-foreground">
              Built with transparency and community in mind
            </p>
          </div>
        </div>
      </section>
    </main>
  )
}