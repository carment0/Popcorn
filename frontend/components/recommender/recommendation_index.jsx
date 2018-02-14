/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import CircularProgress from 'material-ui/CircularProgress';

// Component imports
import MovieTrailerPlayer from './movie_trailer_player';
import YearRangeSelector from './year_range_selector';
import PopularitySelector from './popularity_selector';
import RecommendedMovieItem from './recommended_movie_item';

// Store imports
import { setMovieYearRange } from '../../store/movies/year_range.action';
import { setMoviePopularityPercentile } from '../../store/movies/percentile.action';
import { movieDetailFetch } from '../../store/movies/detail.action';
import { movieTrailerFetch } from '../../store/movies/trailer.action';
import { movieRatingPost, movieRatingRecord } from '../../store/movies/rating.action';
import { movieSkip } from '../../store/movies/movie.action';


// Style imports
import './recommendation_index.scss';


class RecommendationIndex extends React.Component {
  state = {
    trailerSourceList: []
  };

  static propTypes = {
    session: PropTypes.object.isRequired,
    movies: PropTypes.object.isRequired,
    movieTrailers: PropTypes.object.isRequired,
    movieYearRange: PropTypes.object.isRequired,
    movieDetails: PropTypes.object.isRequired,
    moviePopularityPercentile: PropTypes.number.isRequired,
    dispatchSetMovieYearRange: PropTypes.func.isRequired,
    dispatchSetMoviePopularityPercentile: PropTypes.func.isRequired,
    dispatchMovieSkip: PropTypes.func.isRequired,
    dispatchMovieRatingPost: PropTypes.func.isRequired,
    dispatchMovieRatingRecord: PropTypes.func.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired,
    dispatchMovieTrailerFetch: PropTypes.func.isRequired
  }

  playTrailer = (imdbId) => {
    if (this.props.movieTrailers[imdbId] && this.props.movieTrailers[imdbId].length > 0) {
      this.setState({
        trailerSourceList: this.props.movieTrailers[imdbId]
      });
    }
  };

  get filters() {
    return (
      <section className="filters">
        <YearRangeSelector
          disabled={false}
          movieYearRange={this.props.movieYearRange}
          dispatchSetMovieYearRange={this.props.dispatchSetMovieYearRange} />

        <PopularitySelector
          dispatchSetMoviePopularityPercentile={this.props.dispatchSetMoviePopularityPercentile}
          moviePopularityPercentile={this.props.moviePopularityPercentile}
          disabled={false} />
      </section>
    );
  }

  get trailerPlayer() {
    if (this.state.trailerSourceList.length > 0) {
      return (
        <section className="trailer-player">
          <MovieTrailerPlayer videoSourceList={this.state.trailerSourceList} />
        </section>
      );
    }

    const movieIMDBIds = Object.keys(this.props.movieTrailers);
    if (movieIMDBIds.length > 0) {
      const randomIndex = Math.floor(Math.random() * movieIMDBIds.length);
      const randomImdbId = movieIMDBIds[randomIndex];
      return (
        <section className="trailer-player">
          <MovieTrailerPlayer videoSourceList={this.props.movieTrailers[randomImdbId]} />
        </section>
      );
    }

    return <section className="trailer-player" />;
  }

  get recommendedMovies() {
    const validRecommendations = [];

    this.props.movies.recommended.forEach((movieId) => {
      if (this.props.movies.skipped.has(movieId) || this.props.movies.rated.has(movieId)) {
        return;
      }

      if (this.props.movies.all[movieId]) {
        validRecommendations.push(this.props.movies.all[movieId]);
      }
    });

    return validRecommendations.sort().map((movie) => {
      return (
        <RecommendedMovieItem key={movie.id}
          session={this.props.session}
          movieId={movie.id}
          imdbId={movie.imdb_id}
          trailerList={this.props.movieTrailers[movie.imdb_id]}
          movieDetail={this.props.movieDetails[movie.imdb_id]}
          playTrailer={this.playTrailer}
          dispatchMovieSkip={this.props.dispatchMovieSkip}
          dispatchMovieRatingPost={this.props.dispatchMovieRatingPost}
          dispatchMovieRatingRecord={this.props.dispatchMovieRatingRecord}
          dispatchMovieDetailFetch={this.props.dispatchMovieDetailFetch}
          dispatchMovieTrailerFetch={this.props.dispatchMovieTrailerFetch} />
      );
    });
  }

  get movies() {
    if (this.props.movies.recommended.size === 0) {
      return (
        <section className="movies">
          <div className="circular-progress">
            <h4>Please wait while recommendations are fetching...</h4>
            <CircularProgress size={100} thickness={7} />
          </div>
        </section>
      );
    }

    return (
      <section className="movies">
        {this.recommendedMovies}
      </section>
    );
  }

  render() {
    return (
      <section className="recommendation-index">
        <header>
          <h1>Recommendations</h1>
        </header>
        {this.filters}
        {this.trailerPlayer}
        {this.movies}
      </section>
    );
  }
}

const mapReduxStateToProps = (state) => ({
  session: state.session,
  movieDetails: state.movieDetails,
  movieTrailers: state.movieTrailers,
  movieYearRange: state.movieYearRange,
  moviePopularityPercentile: state.moviePopularityPercentile,
  movies: state.movies
});

const mapDispatchToProps = (dispatch) => ({
  dispatchSetMovieYearRange: (minYear, maxYear) => dispatch(setMovieYearRange(minYear, maxYear)),
  dispatchSetMoviePopularityPercentile: (percentile) => dispatch(setMoviePopularityPercentile(percentile)),
  dispatchMovieSkip: (movieId) => dispatch(movieSkip(movieId)),
  dispatchMovieRatingPost: (movieId, rating) => dispatch(movieRatingPost(movieId, rating)),
  dispatchMovieRatingRecord: (movieId, rating) => dispatch(movieRatingRecord(movieId, rating)),
  dispatchMovieDetailFetch: (imdbId) => dispatch(movieDetailFetch(imdbId)),
  dispatchMovieTrailerFetch: (imdbId) => dispatch(movieTrailerFetch(imdbId))
});

export default connect(mapReduxStateToProps, mapDispatchToProps)(RecommendationIndex);
