import React from "react";
import { Tabs, Spin, Radio} from 'antd';
import { AroundMap } from './AroundMap';
import { API_ROOT, GEO_OPTIONS, POS_KEY, AUTH_HEADER, TOKEN_KEY, GOOGLEMAP_API_KEY } from "../constants";
import { Gallery } from "./Gallery";  
import { CreatePostButton } from "./CreatePostButton"

const { TabPane } = Tabs;
const RadioGroup = Radio.Group;



export class Home extends React.Component {
    state =  {
        isLoadingGeolocation: false,
        error: '',
        isLoadingPosts: false,
        posts: [],
        topic: "around"
    }
    componentDidMount() {
        if ("geolocation" in navigator) {
            this.setState({
                isLoadingGeolocation: true
            });
            navigator.geolocation.getCurrentPosition(
                this.onSuccessLoadGeolocation,
                this.onFailedLoadGeolocation,
                GEO_OPTIONS
            );

        }
        else {
            this.setState({
                error: 'Geolocation is not supported.'
            });

        }
    }

    onSuccessLoadGeolocation = (position) => {
            // console.log(position);
            const { latitude, longitude } = position.coords;
            localStorage.setItem(POS_KEY, JSON.stringify({
                lon: longitude,
                lat: latitude
            }));
            this.setState({
                isLoadingGeolocation: false
            });
            this.loadNearbyPosts();
    }
    onFailedLoadGeolocation = () => {
        this.setState({
            isLoadingGeolocation: false,
            error: "Failed to get user location."
        });
    }

    loadNearbyPosts = (center, radius) => {
        this.setState({
            isLoadingPosts: true
        });
        //fire api call
        const {lat, lon} = center? center: JSON.parse(localStorage.getItem(POS_KEY));
        const range = radius? radius: 20;
        const token = localStorage.getItem(TOKEN_KEY);
        fetch(`${API_ROOT}/search?lat=${lat}&lon=${lon}&&range=${range}`, {
            headers: {
                Authorization: `${AUTH_HEADER} ${token}`
            }
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
            throw new Error("Failed to load posts.");
        })
        .then((data) => {
            // console.log(data);
            if (this.state.topic === 'around'){
                this.setState({
                    isLoadingPosts: false, 
                    posts: data? data: [] 
                });
            }

        })
        .catch((e) => {
            this.setState({
                isLoadingPosts: false,
                error: e.message
            });
        })

    }
    getImagePosts = (posts) => {
        const images = posts
            .filter(({type}) => type==='image')
            .map(({ user, url, message }) => ({
                user: user,
                src: url,
                thumbnail: url,
                caption: message,
                thumbnailWidth: 400,
                thumbnailHeight: 300,
            }));
        return <Gallery images={images} />;
    }
    getVideoPosts = (posts) => {
        // using Rol and Col has a problem, when we shrink the browser window, videos will be overlapped.
        const videos = posts
            .filter(({type}) => type==='video')
            .map(({ user, url, message }) => {
                return (
                    <div style={{margin: '10px'}}>
                        <video src={url} controls className="video-block" />
                        <p>{`${user}: ${message}`}</p>
                    </div>
                );
            });
        return (
            <div style={{display: 'flex'}}>
                {videos}
            </div>
        ); 
    }
    getPanelContent = (type) => {
        const { error, isLoadingGeolocation, isLoadingPosts, posts} = this.state;
        if (error){
            return error;
        }
        else if (isLoadingGeolocation){
            return <Spin tip="Loading geo location..." />;
        }
        else if (isLoadingPosts){
            return <Spin tip="Loading posts..." />;
        }
        else if (posts && posts.length > 0){
            return type === 'image'? this.getImagePosts(posts): this.getVideoPosts(posts)
            
        }
        else{
            return "No nearby posts."
        }
    }

    loadFacesAroundTheWorld = () => {
        this.setState({
            isLoadingPosts: true
        });
        //fire api call
        const token = localStorage.getItem(TOKEN_KEY);
        fetch(`${API_ROOT}/cluster?term=face`, {
            headers: {
                Authorization: `${AUTH_HEADER} ${token}`
            }
        })
        .then((response) => {
            if (response.ok) {
                return response.json();
            }
            throw new Error("Failed to load posts.");
        })
        .then((data) => {
            // console.log(data);
            if (this.state.topic === 'face'){
                this.setState({
                    isLoadingPosts: false, 
                    posts: data? data: [] 
                });    
            }
            

        })
        .catch((e) => {
            this.setState({
                isLoadingPosts: false,
                error: e.message
            });
        })

    }

    onTopicChange = (e) => {
        const topic =  e.target.value;
        this.setState({
            topic
        });
        if (topic === 'face'){
            this.loadFacesAroundTheWorld();
        }
        else{
            this.loadNearbyPosts();
        }

    }
    
    render() {
        const operations = <CreatePostButton loadNearbyPosts={this.loadNearbyPosts} />; 
        return (
            <div>
                <RadioGroup onChange={this.onTopicChange} value={this.state.topic}>
                    <Radio value='around'>Posts Around Me</Radio>
                    <Radio value='face'>Faces Around The World</Radio>
                </RadioGroup>
                <Tabs className="main-tabs" tabBarExtraContent={operations}>
                    <TabPane tab="Image Posts" key="1">
                        {this.getPanelContent('image')}
                    </TabPane>
                    <TabPane tab="Video Posts" key="2">
                        {this.getPanelContent('video')}
                    </TabPane>
                    <TabPane tab="Map" key="3">
                        <AroundMap 
                            googleMapURL={`https://maps.googleapis.com/maps/api/js?key=${GOOGLEMAP_API_KEY}&v=3.exp&language=en-US&libraries=geometry,drawing,places`}
                            loadingElement={<div style={{ height: `100%` }} />}
                            containerElement={<div style={{ height: `600px` }} />}
                            mapElement={<div style={{ height: `100%` }} />}
                            posts = {this.state.posts}
                            loadNearbyPosts={this.state.topic === 'around'? this.loadNearbyPosts: this.loadFacesAroundTheWorld}
                        />
                    </TabPane>
                </Tabs>
          </div>
        );
    }
}