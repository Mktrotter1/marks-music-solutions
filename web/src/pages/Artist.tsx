import { useParams } from 'react-router-dom';
import { useArtist, useArtistAlbums } from '../api/hooks';
import { AlbumCard } from '../components/AlbumCard';
import { Users } from 'lucide-react';
import { pluralize } from '../lib/utils';

export function ArtistPage() {
  const { id } = useParams<{ id: string }>();
  const { data: artist, isLoading } = useArtist(id!);
  const { data: albumsData } = useArtistAlbums(id!);

  if (isLoading) {
    return <div className="text-mms-text-secondary">Loading...</div>;
  }

  if (!artist) {
    return <div className="text-mms-text-secondary">Artist not found</div>;
  }

  const albums = albumsData?.items ?? [];

  return (
    <div>
      {/* Artist header */}
      <div className="flex items-end gap-6 mb-8">
        <div className="w-[180px] h-[180px] rounded-full bg-mms-surface flex-shrink-0 flex items-center justify-center overflow-hidden">
          {artist.image_path ? (
            <img
              src={artist.image_path}
              alt={artist.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <Users size={64} className="text-mms-text-tertiary" />
          )}
        </div>

        <div>
          <p className="text-xs uppercase tracking-wider text-mms-text-secondary mb-1">
            Artist
          </p>
          <h1 className="text-4xl font-bold mb-2">{artist.name}</h1>
          <p className="text-sm text-mms-text-secondary">
            {artist.album_count} {pluralize(artist.album_count ?? 0, 'album')} &middot;{' '}
            {artist.track_count} {pluralize(artist.track_count ?? 0, 'track')}
          </p>
        </div>
      </div>

      {/* Discography */}
      {albums.length > 0 && (
        <section>
          <h2 className="text-xl font-bold mb-4">Discography</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {albums.map((album) => (
              <AlbumCard key={album.id} album={album} />
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
