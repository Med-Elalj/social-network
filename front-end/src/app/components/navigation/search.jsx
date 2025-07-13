"use client";
import { useState, useEffect } from "react";
import Styles from "./nav.module.css";
import Image from "next/image";
import { GetData } from "../../../../utils/sendData";
import Link from "next/link";

export function SearchIcon({ onClick, showSearch }) {
  return (
    <div
      className={`${Styles.linkWithIcon} ${showSearch ? Styles.active : ""}`}
      onClick={onClick}
      style={{ cursor: "pointer" }}
    >
      <span className={Styles.iconWrapper}>
        <Image
          src="/search.svg"
          alt="search-hover"
          width={25}
          height={25}
          className={Styles.iconHover}
        />
        <Image
          src="/search.svg"
          alt="search"
          width={25}
          height={25}
          className={Styles.ico}
        />
      </span>
    </div>
  );
}

export function SearchInput({ onClose }) {
  const [query, setQuery] = useState("");

  return (
    <div className={Styles.overlay}>
      <div className={Styles.searchBox}>
        <div className={Styles.header}>
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search..."
            className={Styles.input}
          />
          <button onClick={onClose} className={Styles.closeBtn}>
            ❌
          </button>
        </div>

        {/* Results area */}
        <div className={Styles.results}>
          <ResultList query={query} />
        </div>
      </div>
    </div>
  );
}

export function ResultList({ query }) {
  const [results, setResults] = useState([]);
  const [offset, setOffset] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setOffset(1);
    setResults([]);
  }, [query]);

  useEffect(() => {
    const fetchResults = async () => {
      if (!query) return;
      setLoading(true);
      try {
        const response = await GetData("/api/v1/get/search", {
          query: query,
          offset: offset,
        });
        if (response?.ok) {
          const data = await response.json();
          
          // Fix 1: Properly handle array concatenation
          if (offset === 1) {
            setResults(data.profiles || []);
          } else {
            setResults((prev) => [...prev, ...(data.profiles || [])]);
          }
          
          setHasMore(data.has_more);
        }
      } catch (error) {
        console.error('Error fetching search results:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchResults();
  }, [query, offset]);

  const scrollingHandle = (e) => {
    const element = e.target;
    const scrollBottom =
      element.scrollHeight - element.scrollTop === element.clientHeight;

    if (scrollBottom && hasMore && !loading) {
      setOffset((prev) => prev + 1);
    }
  };

  return (
    <div
      onScroll={scrollingHandle}
      style={{ maxHeight: "300px", overflowY: "auto" }}
    >
      {results?.length > 0 ? (
        results.map((result, index) => (
          <Link
            // Fix 2: Use unique key instead of index
            key={`${result.name}-${index}`}
            className={Styles.resultItem}
            href={
              result.is_group
                ? `/groupes/profile/${result.name}`
                : `/profile/${result.name}`
            }
          >
            <Image
              src={result.pfp?.String ? result.pfp.String : "/iconMale.png"}
              alt="profile"
              width={50}
              height={50}
              style={{ borderRadius: "50%", marginRight: "10px" }}
            />
            <h3>{result.name}</h3>
          </Link>
        ))
      ) : !loading && query ? (
        // Fix 3: Only show "No results" when there's a query
        <h2>No results to display</h2>
      ) : null}

      {loading && (
        <p style={{ padding: "10px", color: "gray" }}>Loading more...</p>
      )}
    </div>
  );
}