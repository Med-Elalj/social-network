"use client";

import { useState, useEffect } from "react";
import Styles from "./nav.module.css";
import Image from "next/image";
import { GetData, SendData } from "@/app/sendData.js";
import Link from "next/link";
import { useNotification } from "@/app/context/NotificationContext";
import { showNotification } from "@/app/utils";

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

export function SearchInput({ onClose, groupId }) {
  const [query, setQuery] = useState("");
  const {showNotification} = useNotification()

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
          <ResultList query={query} groupId={groupId} />
        </div>
      </div>
    </div>
  );
}

export function ResultList({ query, groupId }) {
  const [results, setResults] = useState([]);
  const [offset, setOffset] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [userToSent, setUserToSent] = useState(null)

  useEffect(() => {
    setOffset(1);
    setResults([]);
  }, [query]);

  useEffect(() => {
    const fetchResults = async () => {
      if (!query) return;
      setLoading(true);

      const response = await GetData("/api/v1/get/search", {
        query: query,
        offset: offset,
        groupId: groupId,
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

  useEffect(()=>{
    const requestFetch = async ()=> {
      const response= await SendData('/api/v1/set/sendRequest', userToSent)
      if (response.ok){
        showNotification(`request sent succeffully to ${userToSent.id}`)
      } else {
        showNotification(`error sending request, try again`)
      }
    }
    if (userToSent!==null) {
      requestFetch()
      setUserToSent(null)
    }
  }, [userToSent])

  return (
    <div
      onScroll={scrollingHandle}
      style={{ maxHeight: "300px", overflowY: "auto" }}
    >
      {results?.length > 0 ? (
        results.map((result, index) => (
          <div
            // Fix 2: Use unique key instead of index
            key={`${result.name}-${index}`}
            style={{display:"flex", justifyContent:"space-between", alignItems:"center"}}
            
          >
            <Link
              href={
                result.is_group
                  ? `/groupes/profile/${result.name}`
                  : `/profile/${result.name}`
              }
              className={Styles.resultItem}
              style={{width:"100%"}}
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
            {groupId && <button onClick={()=>{setUserToSent({target:groupId, type:1, receiver_id:result.id})}} style={{width:"150px", height:"30px",backgroundColor:"var(--hover-color)"}}>add user</button>}
          </div>
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
