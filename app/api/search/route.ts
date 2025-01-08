import { NextResponse } from "next/server";
import { exec } from "child_process";
import { promisify } from "util";
import { parseGosearchResults } from "@/lib/parse-results";
import { promises as fs } from "fs";
import path from "path";

const execAsync = promisify(exec);

// Mark this route as dynamic
export const dynamic = 'force-dynamic';

export async function POST(req: Request) {
  try {
    const { query } = await req.json();

    // Define the paths to the temporary data file
    const dataFilePath = path.join('/tmp', 'data.json');

    // Remove the old data file if it exists
    try {
      await fs.unlink(dataFilePath);
    } catch (err) {
      if ((err as NodeJS.ErrnoException).code !== 'ENOENT') {
        throw err;
      }
    }

    // Execute gosearch binary with the query and specify the output file
    const { stdout, stderr } = await execAsync(`/tmp/gosearch ${query} --output=${dataFilePath}`);

    if (stderr) {
      return new NextResponse(stderr, { status: 500 });
    }

    // Read the results from the temporary data file
    const resultsData = await fs.readFile(dataFilePath, 'utf-8');

    // Parse the results
    const results = parseGosearchResults(resultsData);

    return NextResponse.json(results);
  } catch (error) {
    console.error("Search failed:", error);
    return new NextResponse("Internal server error", { status: 500 });
  }
}