/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import Spinner from 'spin';
import Rating from 'react-rating';
import { OverlayTrigger, Popover } from 'react-bootstrap';

// Style
import './movie_item.scss';


class MovieItem extends React.Component {
  static propTypes = {
    movieDetail: PropTypes.object,
    movieId: PropTypes.number.isRequired,
    imdbId: PropTypes.string.isRequired,
    session: PropTypes.object.isRequired,
    dispatchMovieRatingPost: PropTypes.func.isRequired,
    dispatchMovieRatingRecord: PropTypes.func.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired
  };

  static defaultProps = {
    movieDetail: {
      isDefaultProp: true,
      title: 'Loading...',
      year: '....',
      plot: 'Loading...'
    },
  };

  /**
   * Initiate fetching metadata of the current movie and spinning the loading bar so user knows that it is fetching
   * something.
   */
  fetchMovieDetail = () => {
    if (this.props.movieDetail.isDefaultProp) {
      this.spinner = new Spinner().spin();
      document.getElementById(this.props.movieId).appendChild(this.spinner.el);
      this.props.dispatchMovieDetailFetch(this.props.imdbId);
    }
  };

  /**
   * If user is authenticated, rating is saved to database otherwise the rating is cached in redux store and wait until
   * user sign ups and logs in.
   * @param {number} ratingValue
   */
  handleRatingSelect = (value) => {
    let ratingValue = value;
    if (ratingValue <= 0) {
      ratingValue = 0.5;
    }

    if (this.props.session.currentUser === null) {
      this.props.dispatchMovieRatingRecord(this.props.movieId, ratingValue);
    } else {
      this.props.dispatchMovieRatingPost(this.props.movieId, this.props.session.currentUser.id, ratingValue);
    }
  };

  /**
   * @return {React.Element}
   */
  get poster() {
    if (this.props.movieDetail.poster) {
      return (
        <section className="poster">
          <img id={this.props.movieId} src={this.props.movieDetail.poster} alt="poster" />
        </section>
      );
    }

    return <section className="poster" />;
  }

  componentWillReceiveProps(nextProps) {
    if (!nextProps.movieDetail.isDefaultProp && this.spinner) {
      this.spinner.stop();
    }
  }

  componentDidMount() {
    this.fetchMovieDetail();
  }

  render() {
    const titleString = this.props.movieDetail.title + ' - ' + this.props.movieDetail.year;
    const popover = (
      <Popover id="popover-trigger-click-root-close" title={titleString}>
        {this.props.movieDetail.plot}
      </Popover>
    );

    return (
      <article className="default movie-item" id={this.props.movieId}>
        <OverlayTrigger trigger="click" rootClose placement="bottom" overlay={popover}>
          {this.poster}
        </OverlayTrigger>
        <div className="rating-toolbar">
          <div className="star-toolbar-container">
            <Rating fractions={2} onChange={this.handleRatingSelect} />
          </div>
        </div>
      </article>
    );
  }
}

export default MovieItem;
