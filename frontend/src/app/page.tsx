'use client';

import * as React from 'react';
import { useState, useEffect } from 'react';

import Logo from '~/svg/Logo.svg';
import ButtonLink from '@/components/links/ButtonLink';

// In-memory cache to store the fetched data by page and search query
const cache: { [key: string]: any } = {};

export default function HomePage() {
  const [data, setData] = useState<any[]>([]);
  const [sortedData, setSortedData] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isSorting, setIsSorting] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const limit = 10; // Show 10 rows per page

  useEffect(() => {
    fetchData();
  }, [currentPage, searchQuery]);

  // Helper function to generate a cache key based on the page and search query
  const generateCacheKey = (page: number, searchQuery: string) => {
    return `page=${page}&searchQuery=${searchQuery}`;
  };

  const fetchData = async () => {
    const cacheKey = generateCacheKey(currentPage, searchQuery);

    // Check if data for the current page and search query exists in the cache
    if (cache[cacheKey]) {
      console.log(`Data fetched from cache for key: ${cacheKey}`); // Debug logging
      setData(cache[cacheKey].data);
      setTotalPages(cache[cacheKey].totalPages);
      sortData(cache[cacheKey].data);
      return;
    }

    console.log(`Fetching data for page: ${currentPage}, limit: ${limit}`); // Debug logging
    setIsLoading(true);

    try {
      const query = searchQuery ? `&host=${searchQuery}` : '';
      const response = await fetch(`http://localhost:8080/assets?page=${currentPage}&limit=${limit}${query}`);
      const result = await response.json();

      console.log(result); // Debug logging

      // Read the total count from headers
      const totalItems = response.headers.get('X-Total-Count'); // Backend returns this header
      const totalPages = Math.ceil(Number(totalItems) / limit);

      setData(result);
      setTotalPages(totalPages);
      sortData(result);

      // Store the result in the cache
      cache[cacheKey] = { data: result, totalPages };
    } catch (err) {
      console.error('Error fetching data:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const sortData = (fetchedData: any[]) => {
    setIsSorting(true);
    setTimeout(() => {
      const newSortedData = [...fetchedData].sort((a, b) => a.Host.localeCompare(b.Host));
      setSortedData(newSortedData);
      setIsSorting(false);
    }, 1000); // Artificial delay to simulate sorting
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
    setCurrentPage(1); // Reset to first page after search
  };

  const renderedTableRows = sortedData.map((item, index) => (
      <tr key={index} className="border-b">
        <td className="px-4 py-2">{item.ID}</td>
        <td className="px-4 py-2">{item.Host}</td>
        <td className="px-4 py-2">{item.Comment}</td>
        <td className="px-4 py-2">{item.Owner}</td>
        <td className="px-4 py-2">{(item.IPs || []).map((ip: any) => ip.Address).join(', ')}</td>
        <td className="px-4 py-2">{(item.Ports || []).map((port: any) => port.Port).join(', ')}</td>
      </tr>
  ));

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setCurrentPage(newPage);
    }
  };

  return (
      <main>
        <section className="bg-white">
          <div className="layout relative flex min-h-screen flex-col items-center justify-center py-12 text-center">
            <Logo className="w-16" />
            <h1 className="mt-4">Code Challenge</h1>

            <p className="mt-2 text-sm text-gray-800">
              You have complete freedom to present the data here.
            </p>

            <ButtonLink className="mt-6" href="/components" variant="light">
              See all included components
            </ButtonLink>

            {/* Search Box */}
            <div className="mt-4">
              <input
                  type="text"
                  placeholder="Search by host"
                  className="border p-2"
                  value={searchQuery}
                  onChange={handleSearch}
              />
            </div>

            {/* Pagination Controls */}
            <div className="mt-4 flex justify-between w-full max-w-md mx-auto">
              <button
                  disabled={currentPage === 1}
                  onClick={() => handlePageChange(currentPage - 1)}
                  className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
              >
                Previous
              </button>
              <p>
                Page {currentPage} of {totalPages}
              </p>
              <button
                  disabled={currentPage === totalPages}
                  onClick={() => handlePageChange(currentPage + 1)}
                  className="px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
              >
                Next
              </button>
            </div>

            {/* Data Table */}
            <div className="mt-8 w-full max-w-4xl mx-auto bg-gray-100 p-4">
              {isLoading ? (
                  <p>Loading...</p>
              ) : sortedData.length === 0 ? (
                  <p>{isSorting ? 'Sorting...' : 'No data found'}</p>
              ) : (
                  <table className="table-auto w-full text-left bg-white border-collapse border border-gray-300">
                    <thead className="bg-gray-200">
                    <tr>
                      <th className="px-4 py-2">ID</th>
                      <th className="px-4 py-2">Host</th>
                      <th className="px-4 py-2">Comment</th>
                      <th className="px-4 py-2">Owner</th>
                      <th className="px-4 py-2">IPs</th>
                      <th className="px-4 py-2">Ports</th>
                    </tr>
                    </thead>
                    <tbody>
                    {renderedTableRows}
                    </tbody>
                  </table>
              )}
            </div>
          </div>
        </section>
      </main>
  );
}
