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
 * V1alpha1SparkClusterSpecMasterVsphereVolume
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:38:59.380Z[Etc/UTC]")
public class V1alpha1SparkClusterSpecMasterVsphereVolume {
  public static final String SERIALIZED_NAME_FS_TYPE = "fsType";
  @SerializedName(SERIALIZED_NAME_FS_TYPE)
  private String fsType;

  public static final String SERIALIZED_NAME_STORAGE_POLICY_I_D = "storagePolicyID";
  @SerializedName(SERIALIZED_NAME_STORAGE_POLICY_I_D)
  private String storagePolicyID;

  public static final String SERIALIZED_NAME_STORAGE_POLICY_NAME = "storagePolicyName";
  @SerializedName(SERIALIZED_NAME_STORAGE_POLICY_NAME)
  private String storagePolicyName;

  public static final String SERIALIZED_NAME_VOLUME_PATH = "volumePath";
  @SerializedName(SERIALIZED_NAME_VOLUME_PATH)
  private String volumePath;


  public V1alpha1SparkClusterSpecMasterVsphereVolume fsType(String fsType) {
    
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


  public V1alpha1SparkClusterSpecMasterVsphereVolume storagePolicyID(String storagePolicyID) {
    
    this.storagePolicyID = storagePolicyID;
    return this;
  }

   /**
   * Get storagePolicyID
   * @return storagePolicyID
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getStoragePolicyID() {
    return storagePolicyID;
  }


  public void setStoragePolicyID(String storagePolicyID) {
    this.storagePolicyID = storagePolicyID;
  }


  public V1alpha1SparkClusterSpecMasterVsphereVolume storagePolicyName(String storagePolicyName) {
    
    this.storagePolicyName = storagePolicyName;
    return this;
  }

   /**
   * Get storagePolicyName
   * @return storagePolicyName
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getStoragePolicyName() {
    return storagePolicyName;
  }


  public void setStoragePolicyName(String storagePolicyName) {
    this.storagePolicyName = storagePolicyName;
  }


  public V1alpha1SparkClusterSpecMasterVsphereVolume volumePath(String volumePath) {
    
    this.volumePath = volumePath;
    return this;
  }

   /**
   * Get volumePath
   * @return volumePath
  **/
  @ApiModelProperty(required = true, value = "")

  public String getVolumePath() {
    return volumePath;
  }


  public void setVolumePath(String volumePath) {
    this.volumePath = volumePath;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1SparkClusterSpecMasterVsphereVolume v1alpha1SparkClusterSpecMasterVsphereVolume = (V1alpha1SparkClusterSpecMasterVsphereVolume) o;
    return Objects.equals(this.fsType, v1alpha1SparkClusterSpecMasterVsphereVolume.fsType) &&
        Objects.equals(this.storagePolicyID, v1alpha1SparkClusterSpecMasterVsphereVolume.storagePolicyID) &&
        Objects.equals(this.storagePolicyName, v1alpha1SparkClusterSpecMasterVsphereVolume.storagePolicyName) &&
        Objects.equals(this.volumePath, v1alpha1SparkClusterSpecMasterVsphereVolume.volumePath);
  }

  @Override
  public int hashCode() {
    return Objects.hash(fsType, storagePolicyID, storagePolicyName, volumePath);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1SparkClusterSpecMasterVsphereVolume {\n");
    sb.append("    fsType: ").append(toIndentedString(fsType)).append("\n");
    sb.append("    storagePolicyID: ").append(toIndentedString(storagePolicyID)).append("\n");
    sb.append("    storagePolicyName: ").append(toIndentedString(storagePolicyName)).append("\n");
    sb.append("    volumePath: ").append(toIndentedString(volumePath)).append("\n");
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

