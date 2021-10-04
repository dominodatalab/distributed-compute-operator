/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.spark.models;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1SparkClusterSpecMasterPhotonPersistentDisk
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:38:59.380Z[Etc/UTC]")
public class V1alpha1SparkClusterSpecMasterPhotonPersistentDisk {
  public static final String SERIALIZED_NAME_FS_TYPE = "fsType";
  @SerializedName(SERIALIZED_NAME_FS_TYPE)
  private String fsType;

  public static final String SERIALIZED_NAME_PD_I_D = "pdID";
  @SerializedName(SERIALIZED_NAME_PD_I_D)
  private String pdID;


  public V1alpha1SparkClusterSpecMasterPhotonPersistentDisk fsType(String fsType) {
    
    this.fsType = fsType;
    return this;
  }

   /**
   * Get fsType
   * @return fsType
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getFsType() {
    return fsType;
  }


  public void setFsType(String fsType) {
    this.fsType = fsType;
  }


  public V1alpha1SparkClusterSpecMasterPhotonPersistentDisk pdID(String pdID) {
    
    this.pdID = pdID;
    return this;
  }

   /**
   * Get pdID
   * @return pdID
  **/
  @ApiModelProperty(required = true, value = "")

  public String getPdID() {
    return pdID;
  }


  public void setPdID(String pdID) {
    this.pdID = pdID;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1SparkClusterSpecMasterPhotonPersistentDisk v1alpha1SparkClusterSpecMasterPhotonPersistentDisk = (V1alpha1SparkClusterSpecMasterPhotonPersistentDisk) o;
    return Objects.equals(this.fsType, v1alpha1SparkClusterSpecMasterPhotonPersistentDisk.fsType) &&
        Objects.equals(this.pdID, v1alpha1SparkClusterSpecMasterPhotonPersistentDisk.pdID);
  }

  @Override
  public int hashCode() {
    return Objects.hash(fsType, pdID);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1SparkClusterSpecMasterPhotonPersistentDisk {\n");
    sb.append("    fsType: ").append(toIndentedString(fsType)).append("\n");
    sb.append("    pdID: ").append(toIndentedString(pdID)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

