import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useAlbums, useArtists, usePlaylists } from '../api/hooks';
import { AlbumCard } from '../components/AlbumCard';
import { Link } from 'react-router-dom';
import { Users, ListMusic } from 'lucide-react';

type Tab = 'albums' | 'artists' | 'playlists';

export function LibraryPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialTab = (searchParams.get('tab') as Tab) || 'albums';
  const [activeTab, setActiveTab] = useState<Tab>(initialTab);

  const switchTab = (tab: Tab) => {
    setActiveTab(tab);
    setSearchParams({ tab });
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Your Library</h1>

      {/* Tab nav */}
      <div className="flex gap-1 mb-6 border-b border-mms-border">
        {(['albums', 'artists', 'playlists'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => switchTab(tab)}
            className={`px-4 py-2 text-sm font-medium capitalize transition-colors border-b-2 -mb-px ${
              activeTab === tab
                ? 'text-white border-mms-accent'
                : 'text-mms-text-secondary border-transparent hover:text-white'
            }`}
          >
            {tab}
          </button>
        ))}
      </div>

      {/* Tab content */}
      {activeTab === 'albums' && <AlbumsTab />}
      {activeTab === 'artists' && <ArtistsTab />}
      {activeTab === 'playlists' && <PlaylistsTab />}
    </div>
  );
}

function AlbumsTab() {
  const [page, setPage] = useState(0);
  const pageSize = 60;
  const { data, isLoading } = useAlbums(pageSize, page * pageSize);

  if (isLoading) return <p className="text-mms-text-secondary">Loading...</p>;

  const albums = data?.items ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.ceil(total / pageSize);

  return (
    <div>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
        {albums.map((album) => (
          <AlbumCard key={album.id} album={album} />
        ))}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2 mt-8">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="px-3 py-1 text-sm rounded bg-mms-surface text-mms-text-secondary hover:text-white disabled:opacity-30"
          >
            Previous
          </button>
          <span className="text-sm text-mms-text-secondary">
            Page {page + 1} of {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
            disabled={page >= totalPages - 1}
            className="px-3 py-1 text-sm rounded bg-mms-surface text-mms-text-secondary hover:text-white disabled:opacity-30"
          >
            Next
          </button>
        </div>
      )}

      {albums.length === 0 && (
        <p className="text-center text-mms-text-secondary py-10">
          No albums found. Scan your library to get started.
        </p>
      )}
    </div>
  );
}

function ArtistsTab() {
  const [page, setPage] = useState(0);
  const pageSize = 60;
  const { data, isLoading } = useArtists(pageSize, page * pageSize);

  if (isLoading) return <p className="text-mms-text-secondary">Loading...</p>;

  const artists = data?.items ?? [];
  const total = data?.total ?? 0;
  const totalPages = Math.ceil(total / pageSize);

  return (
    <div>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
        {artists.map((artist) => (
          <Link
            key={artist.id}
            to={`/artist/${artist.id}`}
            className="group block rounded-md bg-mms-card p-3 hover:bg-mms-hover transition-colors"
          >
            <div className="aspect-square rounded-full bg-mms-surface mb-3 flex items-center justify-center overflow-hidden">
              {artist.image_path ? (
                <img
                  src={artist.image_path}
                  alt={artist.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <Users size={32} className="text-mms-text-tertiary" />
              )}
            </div>
            <p className="text-sm font-medium truncate text-center">
              {artist.name}
            </p>
            <p className="text-xs text-mms-text-secondary text-center">
              {artist.album_count} albums
            </p>
          </Link>
        ))}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2 mt-8">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="px-3 py-1 text-sm rounded bg-mms-surface text-mms-text-secondary hover:text-white disabled:opacity-30"
          >
            Previous
          </button>
          <span className="text-sm text-mms-text-secondary">
            Page {page + 1} of {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
            disabled={page >= totalPages - 1}
            className="px-3 py-1 text-sm rounded bg-mms-surface text-mms-text-secondary hover:text-white disabled:opacity-30"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}

function PlaylistsTab() {
  const { data, isLoading } = usePlaylists();

  if (isLoading) return <p className="text-mms-text-secondary">Loading...</p>;

  const playlists = data?.items ?? [];

  return (
    <div>
      {playlists.length === 0 ? (
        <div className="text-center py-10">
          <ListMusic size={40} className="mx-auto text-mms-text-tertiary mb-3" />
          <p className="text-mms-text-secondary">No playlists yet</p>
        </div>
      ) : (
        <div className="space-y-1">
          {playlists.map((playlist) => (
            <Link
              key={playlist.id}
              to={`/playlist/${playlist.id}`}
              className="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-mms-hover transition-colors"
            >
              <div className="w-12 h-12 rounded bg-mms-surface flex items-center justify-center">
                <ListMusic size={20} className="text-mms-text-tertiary" />
              </div>
              <div>
                <p className="text-sm font-medium">{playlist.name}</p>
                <p className="text-xs text-mms-text-secondary">
                  {playlist.track_count} tracks
                </p>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
