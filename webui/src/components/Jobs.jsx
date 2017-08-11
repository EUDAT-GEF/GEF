import React from 'react';
import PropTypes from 'prop-types';
import {Row, Col} from 'react-bootstrap';
import {BootstrapTable, TableHeaderColumn} from 'react-bootstrap-table';


const JobStatusIndicator = ({code, tag, message, jobDuration}) => {
    const inProgressColor = {
        color: '#f45d00'
    }
    const errorColor = {
        color: '#ff0000'
    }
    const successColor = {
        color: '#337ab7'
    }
    const progressAnimation = <img src="/images/progress-animation.gif" />;
    let currentProgress;
    let messageColor;

    if (code < 0) {
        currentProgress = progressAnimation;
        messageColor = inProgressColor;
    } else if (code == 0) {
        messageColor = successColor;
    } else {
        messageColor = errorColor;
    }
    return (
        <Row>
            <Col xs={12} sm={3} md={3} style={{fontWeight: 700}}>{tag}</Col>
            <Col xs={12} sm={9} md={9} style={messageColor}>{currentProgress} {message} (elapsed time {jobDuration})</Col>
        </Row>
    )
};



const selectRowProp = {
    mode: 'checkbox'
};


//let order = 'desc';
let allJobs = [];
let jobStatusUpdateTimer;

class Jobs extends React.Component {
// class MultiSelectTable extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            showModal: false,
            buttonPressed: -1,
            timerOn: true,
        };

        this.options = {
            defaultSortName: 'title',  // default sort column name
            defaultSortOrder: 'desc'  // default sort order
        };

        // if (this.hasJobsRunning() == false) {
        //     this.setState({timerOn: false});
        //     clearInterval(stateUpdateTimer);
        //     console.log("NOT RUNNING");
        // } else {
        jobStatusUpdateTimer = setInterval(this.tick.bind(this), 1000);

        console.log(this.state.timerOn);
    }

    componentDidMount() {
        this.props.fetchJobs();
        this.props.fetchServices();
    }

    componentWillUnmount() {
        clearInterval(jobStatusUpdateTimer);
    }

    tick() {
        this.props.fetchJobs();
        console.log(this.state.timerOn);
        console.log(this.hasJobsRunning());
        if ((this.state.timerOn) && (!this.hasJobsRunning())) {

            clearInterval(jobStatusUpdateTimer);
            this.setState({timerOn: false});
            console.log("NOT RUNNING");
        }

        if ((!this.state.timerOn) && (this.hasJobsRunning())) {
            jobStatusUpdateTimer = setInterval(this.tick.bind(this), 1000);
            this.setState({timerOn: true});
            console.log("----------- RUNNING");
        }
        //this.props.consoleOutputFetch(this.props.job.ID);
    }

    formatJobDuration(durationTime) {
        var sec_num = parseInt(durationTime, 10);
        var hours   = Math.floor(sec_num / 3600);
        var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
        var seconds = sec_num - (hours * 3600) - (minutes * 60);

        if (hours   < 10) {hours   = "0"+hours;}
        if (minutes < 10) {minutes = "0"+minutes;}
        if (seconds < 10) {seconds = "0"+seconds;}
        return hours+':'+minutes+':'+seconds;
    }


    hasJobsRunning() {
        var runningJobfound = false
        this.props.jobs.map((job) => {

            if (job.State.Code == -1) {
                console.log(job.State);
                runningJobfound = true;

            }
        });
        return runningJobfound;
    }




    render() {

        allJobs = [];
        if (this.props.jobs) {
            return (
                <div>
                    <h3>Browse Jobs</h3>
                    <h4>All jobs</h4>
                    {this.props.jobs.map((job) => {
                        let service = null;
                        for (var i = 0; i < this.props.services.length; ++i) {
                            if (job.ServiceID == this.props.services[i].ID) {
                                service = this.props.services[i];
                                break;
                            }
                        }
                        let serviceName = (service && service.Name && service.Name.length) ? service.Name :
                            (service && service.ID && service.ID.length) ? service.ID : "unknown service";
                        let title = "Job from " + serviceName;



                        // if (this.hasJobsRunning() == true) {
                        //     console.log("RUNNING");
                        //     if (this.timerOn == false) {
                        //         this.setState({timerOn: true});
                        //         stateUpdateTimer = setInterval(this.tick.bind(this), 1000);
                        //     }
                        // }


                        // } else {
                        //     this.setState({timerOn: false});
                        //     clearInterval(stateUpdateTimer);
                        //     console.log("NOT RUNNING");
                        // }


                        // if (!this.state.progressIndicator) {
                        //     this.state.progressIndicator = " ";
                        // } else {
                        //     if (this.state.progressIndicator.length>4) {
                        //         this.state.progressIndicator = " ";
                        //     } else {
                        //         this.state.progressIndicator += ".";
                        //     }
                        // }

                        let execDuration = "";
                        if (job.State.Code == -1) {
                            // clearInterval(stateUpdateTimer);
                            // buttonClass = "btn btn-default";
                            // this.state.progressIndicator = "";

                            let currentDate = new Date();
                            execDuration = currentDate - Date.parse(job.Created);
                        } else {
                            execDuration = Date.parse(job.Finished) - Date.parse(job.Created);
                        }







                        let createdDate = new Date(job.Created);

                        let fmtCreatedDate = createdDate.toLocaleDateString('en-GB');
                        let fmtCreatedTime = createdDate.toLocaleTimeString('en-GB');


                        allJobs.push({"title": title, "id": job.ID, "created": fmtCreatedDate + " " + fmtCreatedTime, "duration": this.formatJobDuration(execDuration/1000), "status": job.State.Status});
                    })}

                    <div>
                        <BootstrapTable data={allJobs} selectRow={selectRowProp} options={this.options}>
                            <TableHeaderColumn dataField='id' isKey dataSort>ID</TableHeaderColumn>
                            <TableHeaderColumn dataField='title' dataSort>Title</TableHeaderColumn>
                            <TableHeaderColumn dataField='created' dataSort>Created</TableHeaderColumn>
                            <TableHeaderColumn dataField='duration' dataSort>Duration</TableHeaderColumn>
                            <TableHeaderColumn dataField='status' dataSort>Status</TableHeaderColumn>
                        </BootstrapTable>
                    </div>
                </div>
            );
        } else {
            return (
                <div><h4>No jobs found</h4></div>
            )
        }










    //
    //
    //     if (this.props.jobs) {
    //         return (
    //             <h3>Browse Jobs</h3>
    //         { this.props.jobs.map((job) => {
    //             console.log("TEXT");
    //             console.log(job);
    //             allJobs.push({"id": job.ID, "serviceID": job.ServiceID, "status": job.State.Status})
    //
    //
    //         })
    //         }
    //         // console.log(allJobs);
    //         // console.log(products);
    //
    //
    //         <div>
    //             <BootstrapTable data={products} selectRow={selectRowProp} options={this.options}>
    //                 <TableHeaderColumn dataField='id' isKey dataSort>Product ID</TableHeaderColumn>
    //                 <TableHeaderColumn dataField='name' dataSort>Product Name</TableHeaderColumn>
    //                 <TableHeaderColumn dataField='price'>Product Price</TableHeaderColumn>
    //             </BootstrapTable>
    //         </div>
    //     )
    //     }
    // }
    //
    //
    //




    }

}

Jobs.propTypes = {
//MultiSelectTable.propTypes = {
    jobs: PropTypes.array, // can be null
    fetchJobs: PropTypes.func.isRequired,
    services: PropTypes.array, // can be null
    fetchServices: PropTypes.func.isRequired,
    jobID: PropTypes.string
};

export default Jobs;
//export default MultiSelectTable;