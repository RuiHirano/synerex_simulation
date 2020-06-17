
import React, { useEffect, useState } from 'react';
import { HarmoVisLayers, Container, BasedProps, BasedState, connectToHarmowareVis, MovesLayer, Movesbase, MovesbaseOperation, DepotsLayer, DepotsData, LineMapLayer, LineMapData } from 'harmoware-vis';
import io from "socket.io-client";
import { Controller } from '../components';


//const MAPBOX_TOKEN = process.env.REACT_APP_MAPBOX_TOKEN ? process.env.REACT_APP_MAPBOX_TOKEN : "";
const MAPBOX_TOKEN = 'pk.eyJ1IjoicnVpaGlyYW5vIiwiYSI6ImNqdmc0bXJ0dTAzZDYzem5vMmk0ejQ0engifQ.3k045idIb4JNvawjppzqZA'


class Harmoware extends Container<BasedProps & BasedState> {
    render() {
        const { actions, depotsData, viewport } = this.props;
        return (<HarmowarePage {...this.props} />)
    }
}

const socket: SocketIOClient.Socket = io();

interface AreaInfo {
    Id: string
    Name: string
    ControlArea: Coord[]
    DuplicateArea: Coord[]
}

interface Coord {
    Lat: number
    Lon: number
}

const HarmowarePage: React.FC<BasedProps> = (props) => {
    const { actions, depotsData, viewport, movesbase, movedData, routePaths, clickedObject } = props

    const [linedata, setLinedata] = useState<LineMapData[]>([])
    const [areadata, setAreadata] = useState<AreaInfo[]>([])
    const [movesbases, setMovesbases] = useState<AreaInfo[]>([])

    const getAgents = (data: any) => {
        const time = Date.now() / 1000; // set time as now. (If data have time, ..)
        console.log("socketData length2", data.length);
        //console.log("movesbasedata length", movesbasedata.length)
        const movesbases: Movesbase[] = [];

        data.forEach((value: any) => {
            const { mtype, id, lat, lon } = JSON.parse(
                value
            );

            let color = [0, 200, 120];
            movesbases.push({
                type: mtype,
                movesbaseidx: id,
                departuretime: time,
                arrivaltime: time,
                operation: [
                    {
                        elapsedtime: time,
                        position: [lon, lat, 0],
                        //direction: 10,
                        color
                    }
                ]
            });
        });

        actions.updateMovesBase(movesbases);
    }

    const getAreas = (data: any) => {
        console.log("areaInfo", data);

        const linedata: LineMapData[] = []
        const areas = convertJsonToArea(data)
        setAreadata(areas)

        areas.forEach((areaInfo: AreaInfo) => {
            const { maxLat, maxLon, minLat, minLon } = getCoordRange(areaInfo.ControlArea)
            linedata.push({
                "sourcePosition": [minLon, minLat, 0],
                "targetPosition": [minLon, maxLat, 0]
            })
            linedata.push({
                "sourcePosition": [minLon, maxLat, 0],
                "targetPosition": [maxLon, maxLat, 0]
            })
            linedata.push({
                "sourcePosition": [maxLon, maxLat, 0],
                "targetPosition": [maxLon, minLat, 0]
            })
            linedata.push({
                "sourcePosition": [maxLon, minLat, 0],
                "targetPosition": [minLon, minLat, 0]
            })
        })

        setLinedata(linedata)
    }

    /*const testAgent = async () => {
        for (let i = 0; i < 100; i++) {
            const setMovesbase: Movesbase[] = [];
            const time = Date.now() / 1000;
            let color = [0, 200, 0];
            await timeout(1000)
            for (let index = 0; index < 100; index++) {
                setMovesbase.push({
                    type: 'ped',
                    movesbaseidx: index,
                    departuretime: time,
                    arrivaltime: time,
                    operation: [
                        {
                            elapsedtime: time,
                            position: [135.4664 + index * 0.0001, 35.633253 + index * 0.0001, 0],
                            direction: 10,
                            color
                        }
                    ]
                });
            }

            actions.updateMovesBase(setMovesbase);
        }
    }*/

    useEffect(() => {
        socket.on("agents", (data: any) => getAgents(data));
        socket.on("areas", (data: any) => getAreas(data));

        //testAgent()
        console.log(process.env);
        if (actions) {
            actions.setViewport({
                ...props.viewport,
                width: window.screen.width,
                height: window.screen.height,
                zoom: 10
            })
            actions.setSecPerHour(1000);
        }
    }, [])
    return (
        <div>
            <Controller {...props} />
            <HarmoVisLayers
                viewport={viewport} actions={actions}
                mapboxApiAccessToken={MAPBOX_TOKEN}
                layers={[
                    new LineMapLayer({
                        data: linedata,
                        getWidth: (x) => 20,
                    }),
                    new DepotsLayer({
                        depotsData,
                        iconChange: false,
                        layerRadiusScale: 20
                    }),
                    new MovesLayer({
                        routePaths,
                        movesbase,
                        movedData,
                        clickedObject,
                        actions,
                        //lightSettings,
                        layerRadiusScale: 0.1,
                        getRadius: x => 1,
                        getRouteWidth: x => 1,
                        optionCellSize: 2,
                        sizeScale: 1,
                        iconChange: false,
                        optionChange: false, // this.state.optionChange,
                        //onHover
                    })
                ]}
            />
        </div>
    );
}

/*async function timeout(ms: number) {
    await new Promise(resolve => setTimeout(resolve, ms));
    return
}*/

const getCoordRange = ((coords: Coord[]) => {
    let maxLat = Number.NEGATIVE_INFINITY
    let maxLon = Number.NEGATIVE_INFINITY
    let minLat = Number.POSITIVE_INFINITY
    let minLon = Number.POSITIVE_INFINITY

    coords.forEach((coord) => {
        if (coord.Lat > maxLat) {
            maxLat = coord.Lat
        }
        if (coord.Lon > maxLon) {
            maxLon = coord.Lon
        }
        if (coord.Lat < minLat) {
            minLat = coord.Lat
        }
        if (coord.Lon < minLon) {
            minLon = coord.Lon
        }
    })

    return { maxLat, maxLon, minLat, minLon }
})

const convertJsonToArea = ((data: any[]) => {
    const areas: AreaInfo[] = []
    data.forEach((areaStr: any) => {
        const areaJson = JSON.parse(areaStr);
        var area: AreaInfo = { ControlArea: [], DuplicateArea: [], Name: "", Id: "" }
        areaJson.control_area.forEach((arg: any) => {
            area.ControlArea.push({ Lat: arg.latitude, Lon: arg.longitude })
        })
        areaJson.duplicate_area.forEach((arg: any) => {
            area.DuplicateArea.push({ Lat: arg.latitude, Lon: arg.longitude })
        })
        areaJson.id ? area.Id = areaJson.id : area.Id = ""
        areaJson.name ? area.Name = areaJson.name : area.Name = ""
        areas.push(area)
    })

    return areas
})

export default connectToHarmowareVis(Harmoware);